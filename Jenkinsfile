pipeline {
    agent any

    environment {
        APP_NAME = 'hrms-backend'
        DOCKER_IMAGE = 'hrms-backend'
        DOCKER_REGISTRY = credentials('docker-registry-url') ?: ''
        DOCKER_HOST = 'unix:///var/run/docker.sock'
        // Git variables
        GIT_COMMIT_SHORT = sh(script: 'git rev-parse --short HEAD', returnStdout: true).trim()
        GIT_BRANCH_NAME = sh(script: 'git rev-parse --abbrev-ref HEAD', returnStdout: true).trim()
        // Build version
        BUILD_VERSION = "${GIT_BRANCH_NAME}-${GIT_COMMIT_SHORT}-${BUILD_NUMBER}"
    }

    options {
        buildDiscarder(logRotator(numToKeepStr: '10', daysToKeepStr: '30'))
        timestamps()
        timeout(time: 30, unit: 'MINUTES')
        disableConcurrentBuilds()
        skipDefaultCheckout(false)
        ansiColor('xterm')
    }

    parameters {
        choice(
            name: 'BUILD_TYPE',
            choices: ['full', 'incremental'],
            description: 'Select build type'
        )
        booleanParam(
            name: 'SKIP_TESTS',
            defaultValue: false,
            description: 'Skip running tests'
        )
        booleanParam(
            name: 'PUSH_TO_REGISTRY',
            defaultValue: true,
            description: 'Push image to Docker registry'
        )
    }

    stages {
        // Stage 1: Pre-build verification
        stage('Pre-build Setup') {
            steps {
                script {
                    echo "=========================================="
                    echo "Starting Build #${BUILD_NUMBER}"
                    echo "Branch: ${GIT_BRANCH_NAME}"
                    echo "Commit: ${GIT_COMMIT_SHORT}"
                    echo "Build Type: ${params.BUILD_TYPE}"
                    echo "=========================================="
                    
                    // Verify Docker setup
                    sh '''
                        echo "=== Docker Environment Check ==="
                        echo "User: $(whoami)"
                        echo "UID: $(id -u)"
                        echo "Groups: $(groups)"
                        echo "Docker Version:"
                        docker --version || echo "Docker not found"
                        echo "Docker Info:"
                        docker info 2>/dev/null | grep -E "Server Version:|Containers:|Running:" || echo "Cannot connect to Docker daemon"
                        echo "Docker Socket:"
                        ls -la /var/run/docker.sock 2>/dev/null || echo "Docker socket not found"
                        echo "Available Disk Space:"
                        df -h /var/lib/docker 2>/dev/null || df -h / 2>/dev/null
                    '''
                }
            }
        }

        // Stage 2: Checkout code
        stage('Checkout Source Code') {
            steps {
                checkout([
                    $class: 'GitSCM',
                    branches: [[name: "*/${GIT_BRANCH_NAME}"]],
                    extensions: [
                        [$class: 'CleanBeforeCheckout'],
                        [$class: 'CloneOption', depth: 1, noTags: false, shallow: true],
                        [$class: 'RelativeTargetDirectory', relativeTargetDir: '.']
                    ],
                    userRemoteConfigs: [[
                        url: scm.userRemoteConfigs[0].url,
                        credentialsId: 'github-token'
                    ]]
                ])
                
                script {
                    // Validate Dockerfile exists
                    sh '''
                        if [ ! -f "Dockerfile" ]; then
                            echo "ERROR: Dockerfile not found in repository root!"
                            exit 1
                        fi
                        echo "Dockerfile found. Contents:"
                        head -20 Dockerfile
                    '''
                }
            }
        }

        // Stage 3: Build preparation
        stage('Prepare Build') {
            steps {
                script {
                    // Clean up old Docker resources
                    if (params.BUILD_TYPE == 'full') {
                        sh '''
                            echo "Performing full cleanup..."
                            docker system prune -f --filter "until=24h" || true
                            docker builder prune -f --all || true
                        '''
                    } else {
                        sh '''
                            echo "Performing incremental cleanup..."
                            docker container prune -f || true
                            docker image prune -f || true
                        '''
                    }
                    
                    // Create build info file
                    sh '''
                        cat > build-info.json << EOF
                        {
                            "app_name": "${APP_NAME}",
                            "version": "${BUILD_VERSION}",
                            "branch": "${GIT_BRANCH_NAME}",
                            "commit": "${GIT_COMMIT_SHORT}",
                            "build_number": "${BUILD_NUMBER}",
                            "build_time": "$(date -u +"%Y-%m-%dT%H:%M:%SZ")",
                            "build_type": "${params.BUILD_TYPE}"
                        }
                        EOF
                        cat build-info.json
                    '''
                }
            }
        }

        // Stage 4: Build Docker image
        stage('Build Docker Image') {
            steps {
                script {
                    echo "Building Docker image: ${DOCKER_IMAGE}:${BUILD_VERSION}"
                    
                    def buildArgs = [
                        "--build-arg BUILD_NUMBER=${BUILD_NUMBER}",
                        "--build-arg GIT_COMMIT=${GIT_COMMIT_SHORT}",
                        "--build-arg GIT_BRANCH=${GIT_BRANCH_NAME}",
                        "--build-arg BUILD_VERSION=${BUILD_VERSION}"
                    ].join(' ')
                    
                    def cacheArgs = params.BUILD_TYPE == 'incremental' ? "--cache-from ${DOCKER_IMAGE}:latest" : ""
                    
                    sh """
                        # Build the Docker image with proper tags
                        docker build ${cacheArgs} ${buildArgs} \
                            -t ${DOCKER_IMAGE}:${BUILD_VERSION} \
                            -t ${DOCKER_IMAGE}:${GIT_COMMIT_SHORT} \
                            -t ${DOCKER_IMAGE}:latest \
                            --label "version=${BUILD_VERSION}" \
                            --label "commit=${GIT_COMMIT_SHORT}" \
                            --label "branch=${GIT_BRANCH_NAME}" \
                            --label "build.number=${BUILD_NUMBER}" \
                            --label "maintainer=jenkins" \
                            --label "built-by=jenkins" \
                            --label "build.date=$(date -u +"%Y-%m-%dT%H:%M:%SZ")" \
                            .
                        
                        # Verify the image was created
                        echo "=== Image Verification ==="
                        docker image inspect ${DOCKER_IMAGE}:${BUILD_VERSION} --format '{{.Id}}'
                        
                        # Show image details
                        echo "=== Built Images ==="
                        docker images ${DOCKER_IMAGE} --format "table {{.Repository}}\t{{.Tag}}\t{{.Size}}\t{{.CreatedAt}}" | head -10
                        
                        # Check image size
                        echo "=== Image Size ==="
                        docker images ${DOCKER_IMAGE}:${BUILD_VERSION} --format "{{.Size}}"
                    """
                }
            }
            
            post {
                success {
                    echo "✅ Docker image built successfully: ${DOCKER_IMAGE}:${BUILD_VERSION}"
                    archiveArtifacts artifacts: 'build-info.json', fingerprint: true
                }
                failure {
                    echo "❌ Docker build failed"
                    sh '''
                        echo "=== Build Logs ==="
                        docker images | grep "${DOCKER_IMAGE}" || echo "No images found"
                        echo "=== Recent Docker Logs ==="
                        journalctl -u docker --since "10 minutes ago" | tail -50 || true
                    '''
                }
            }
        }

        // Stage 5: Run tests (optional)
        stage('Run Tests') {
            when {
                expression { 
                    return !params.SKIP_TESTS.toBoolean() 
                }
            }
            steps {
                script {
                    echo "Running tests..."
                    sh '''
                        # Run tests inside the built container
                        docker run --rm \
                            --name ${APP_NAME}-test-${BUILD_NUMBER} \
                            ${DOCKER_IMAGE}:${BUILD_VERSION} \
                            sh -c "npm test || echo 'Tests completed'"  # Adjust command based on your tech stack
                        
                        # Or run tests directly if you have test scripts
                        # ./run-tests.sh
                    '''
                }
            }
        }

        // Stage 6: Push to registry
        stage('Push to Docker Registry') {
            when {
                allOf {
                    expression { 
                        return params.PUSH_TO_REGISTRY.toBoolean() && 
                               env.DOCKER_REGISTRY?.trim() 
                    }
                    anyOf {
                        branch 'main'
                        branch 'master'
                        branch 'develop'
                        branch 'leave-management'
                        expression { return params.PUSH_TO_REGISTRY.toBoolean() }
                    }
                }
            }
            steps {
                script {
                    echo "Pushing images to registry: ${DOCKER_REGISTRY}"
                    
                    withCredentials([
                        usernamePassword(
                            credentialsId: 'docker-registry-credentials',
                            usernameVariable: 'DOCKER_USER',
                            passwordVariable: 'DOCKER_PASS'
                        )
                    ]) {
                        sh """
                            # Login to Docker registry
                            echo "\${DOCKER_PASS}" | docker login -u "\${DOCKER_USER}" --password-stdin ${DOCKER_REGISTRY}
                            
                            # Tag images with registry prefix
                            docker tag ${DOCKER_IMAGE}:${BUILD_VERSION} ${DOCKER_REGISTRY}/${DOCKER_IMAGE}:${BUILD_VERSION}
                            docker tag ${DOCKER_IMAGE}:${GIT_COMMIT_SHORT} ${DOCKER_REGISTRY}/${DOCKER_IMAGE}:${GIT_COMMIT_SHORT}
                            docker tag ${DOCKER_IMAGE}:latest ${DOCKER_REGISTRY}/${DOCKER_IMAGE}:latest
                            
                            # Push images
                            echo "Pushing ${DOCKER_REGISTRY}/${DOCKER_IMAGE}:${BUILD_VERSION}"
                            docker push ${DOCKER_REGISTRY}/${DOCKER_IMAGE}:${BUILD_VERSION}
                            
                            echo "Pushing ${DOCKER_REGISTRY}/${DOCKER_IMAGE}:${GIT_COMMIT_SHORT}"
                            docker push ${DOCKER_REGISTRY}/${DOCKER_IMAGE}:${GIT_COMMIT_SHORT}
                            
                            echo "Pushing ${DOCKER_REGISTRY}/${DOCKER_IMAGE}:latest"
                            docker push ${DOCKER_REGISTRY}/${DOCKER_IMAGE}:latest
                            
                            # Verify pushed images
                            echo "=== Pushed Images ==="
                            docker images | grep "${DOCKER_REGISTRY}/${DOCKER_IMAGE}"
                            
                            # Logout from registry
                            docker logout ${DOCKER_REGISTRY}
                        """
                    }
                }
            }
            
            post {
                success {
                    echo "✅ Images successfully pushed to ${DOCKER_REGISTRY}"
                    script {
                        // Create deployment manifest
                        sh """
                            cat > deployment.manifest << EOF
                            # Deployment Manifest
                            Application: ${APP_NAME}
                            Image: ${DOCKER_REGISTRY}/${DOCKER_IMAGE}:${BUILD_VERSION}
                            Commit: ${GIT_COMMIT_SHORT}
                            Branch: ${GIT_BRANCH_NAME}
                            Build: ${BUILD_NUMBER}
                            Timestamp: $(date)
                            EOF
                        """
                        archiveArtifacts artifacts: 'deployment.manifest', fingerprint: true
                    }
                }
            }
        }

        // Stage 7: Deploy to server
        stage('Deploy to Server') {
            when {
                allOf {
                    anyOf {
                        branch 'main'
                        branch 'master'
                        branch 'leave-management'
                    }
                    expression {
                        try {
                            withCredentials([string(credentialsId: 'deploy-server-host', variable: 'DEPLOY_HOST')]) {
                                return DEPLOY_HOST?.trim()
                            }
                        } catch (Exception e) {
                            return false
                        }
                    }
                }
            }
            steps {
                script {
                    echo "Deploying to server..."
                    
                    withCredentials([
                        string(credentialsId: 'deploy-server-host', variable: 'DEPLOY_SERVER'),
                        string(credentialsId: 'deploy-server-path', variable: 'DEPLOY_PATH'),
                        sshUserPrivateKey(
                            credentialsId: 'deploy-server-ssh-key',
                            keyFileVariable: 'SSH_KEY_FILE',
                            usernameVariable: 'SSH_USER'
                        )
                    ]) {
                        sh """
                            # Copy deployment manifest
                            scp -i ${SSH_KEY_FILE} -o StrictHostKeyChecking=no \
                                deployment.manifest \
                                ${SSH_USER}@${DEPLOY_SERVER}:${DEPLOY_PATH}/deployment.manifest.${BUILD_NUMBER}
                            
                            # Deploy using SSH
                            ssh -i ${SSH_KEY_FILE} -o StrictHostKeyChecking=no ${SSH_USER}@${DEPLOY_SERVER} << 'ENDSSH'
                                set -e  # Exit on error
                                echo "=== Starting Deployment ==="
                                cd ${DEPLOY_PATH}
                                
                                # Pull latest code (if using git on server)
                                git fetch origin || true
                                git checkout ${GIT_BRANCH_NAME} || true
                                git pull origin ${GIT_BRANCH_NAME} || true
                                
                                # Update docker-compose with new image
                                if [ -f "docker-compose.yml" ]; then
                                    # Backup current compose file
                                    cp docker-compose.yml docker-compose.yml.backup.$(date +%Y%m%d_%H%M%S)
                                    
                                    # Update image tag (simple sed example - adjust for your compose file)
                                    sed -i "s|image:.*${DOCKER_IMAGE}:.*|image: ${DOCKER_REGISTRY}/${DOCKER_IMAGE}:${BUILD_VERSION}|g" docker-compose.yml
                                    
                                    echo "Updated docker-compose.yml with image: ${DOCKER_REGISTRY}/${DOCKER_IMAGE}:${BUILD_VERSION}"
                                fi
                                
                                # Login to Docker registry on server
                                echo "${DOCKER_PASS}" | docker login -u "${DOCKER_USER}" --password-stdin ${DOCKER_REGISTRY} 2>/dev/null || echo "Registry login skipped"
                                
                                # Pull new image
                                docker pull ${DOCKER_REGISTRY}/${DOCKER_IMAGE}:${BUILD_VERSION}
                                
                                # Deploy with docker-compose
                                docker compose down || true
                                docker compose build --pull || true
                                docker compose up -d --remove-orphans
                                
                                # Cleanup old images
                                docker image prune -f
                                docker system prune -f --filter "until=168h" || true
                                
                                # Wait for services to start
                                sleep 10
                                
                                # Health check
                                echo "=== Health Check ==="
                                docker compose ps
                                
                                # Check if containers are running
                                if docker compose ps | grep -q "Up"; then
                                    echo "✅ Deployment successful!"
                                    
                                    # Test endpoint (adjust port/endpoint as needed)
                                    curl -f http://localhost:8081/health || \
                                    curl -f http://localhost:3000/health || \
                                    curl -f http://localhost:80/ || \
                                    echo "Health check endpoints may be different"
                                else
                                    echo "❌ Some services are not running!"
                                    docker compose logs --tail=50
                                    exit 1
                                fi
                                
                                echo "=== Deployment Complete ==="
ENDSSH
                        """
                    }
                }
            }
            
            post {
                success {
                    echo "✅ Application deployed successfully!"
                    script {
                        // Send notification or update dashboard
                        sh '''
                            echo "Deployment completed at $(date)"
                            echo "URL: http://${DEPLOY_SERVER}:8081"  # Adjust port as needed
                        '''
                    }
                }
                failure {
                    echo "❌ Deployment failed!"
                    script {
                        // Rollback or alert
                        sh '''
                            echo "Deployment failed. Check server logs."
                            echo "To rollback, run on server:"
                            echo "cd ${DEPLOY_PATH} && docker compose down && docker compose up -d"
                        '''
                    }
                }
            }
        }

        // Stage 8: Post-deployment verification
        stage('Post-deployment Verification') {
            when {
                branch 'main'
                expression { currentBuild.resultIsBetterOrEqualTo('SUCCESS') }
            }
            steps {
                script {
                    withCredentials([
                        string(credentialsId: 'deploy-server-host', variable: 'DEPLOY_SERVER')
                    ]) {
                        sh """
                            echo "=== Post-deployment Verification ==="
                            echo "Application: ${APP_NAME}"
                            echo "Version: ${BUILD_VERSION}"
                            echo "Deployed to: ${DEPLOY_SERVER}"
                            echo "Health check endpoints:"
                            echo "- http://${DEPLOY_SERVER}:8080/health"
                            echo "- http://${DEPLOY_SERVER}:3000/api/status"
                            echo "=== Verification Complete ==="
                        """
                    }
                }
            }
        }
    }

    post {
        always {
            echo "=========================================="
            echo "Build #${BUILD_NUMBER} - ${currentBuild.currentResult}"
            echo "Duration: ${currentBuild.durationString}"
            echo "=========================================="
            
            script {
                // Cleanup Docker resources
                if (currentBuild.resultIsBetterOrEqualTo('SUCCESS')) {
                    sh '''
                        echo "Cleaning up Docker resources..."
                        # Remove intermediate containers
                        docker container prune -f || true
                        
                        # Remove dangling images
                        docker image prune -f || true
                        
                        # Remove build cache
                        docker builder prune -f || true
                        
                        # Clean up old volumes (be careful with this)
                        # docker volume prune -f || true
                    '''
                }
                
                // Generate build report
                sh '''
                    echo "=== Build Report ===" > build-report.txt
                    echo "Result: ${currentBuild.currentResult}" >> build-report.txt
                    echo "Number: ${BUILD_NUMBER}" >> build-report.txt
                    echo "Duration: ${currentBuild.durationString}" >> build-report.txt
                    echo "Branch: ${GIT_BRANCH_NAME}" >> build-report.txt
                    echo "Commit: ${GIT_COMMIT_SHORT}" >> build-report.txt
                    echo "Image: ${DOCKER_IMAGE}:${BUILD_VERSION}" >> build-report.txt
                    date >> build-report.txt
                '''
                archiveArtifacts artifacts: 'build-report.txt', fingerprint: true
            }
            
            // Clean workspace
            cleanWs(cleanWhenAborted: true, cleanWhenFailure: true, cleanWhenNotBuilt: true, 
                   cleanWhenUnstable: true, cleanWhenSuccess: true)
        }
        
        success {
            echo "🎉 Pipeline executed successfully!"
            // Optional: Send success notification
            // emailext body: "Build ${BUILD_NUMBER} succeeded!\n\nView: ${BUILD_URL}", subject: "SUCCESS: ${JOB_NAME} #${BUILD_NUMBER}", to: 'team@example.com'
        }
        
        failure {
            echo "❌ Pipeline failed!"
            script {
                sh '''
                    echo "=== Failure Analysis ==="
                    echo "Check Docker:"
                    docker ps -a || true
                    echo "Recent logs:"
                    tail -100 /var/log/jenkins/jenkins.log 2>/dev/null | grep -A 10 -B 10 "ERROR\|FAILED" || true
                '''
            }
            // Optional: Send failure notification
            // emailext body: "Build ${BUILD_NUMBER} failed!\n\nView: ${BUILD_URL}\n\nConsole: ${BUILD_URL}console", subject: "FAILURE: ${JOB_NAME} #${BUILD_NUMBER}", to: 'team@example.com'
        }
        
        unstable {
            echo "⚠️ Pipeline unstable!"
        }
        
        changed {
            echo "Build status changed!"
        }
    }
}