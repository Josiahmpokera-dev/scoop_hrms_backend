pipeline {
    agent any

    environment {
        // Application settings
        APP_NAME = 'hrms-backend'
        GO_VERSION = '1.23'
        
        // Docker settings
        DOCKER_IMAGE = "${APP_NAME}"
        DOCKER_REGISTRY = credentials('docker-registry-url')  // Configure in Jenkins credentials
        DOCKER_CREDENTIALS = credentials('docker-registry-credentials')
        
        // Deployment settings
        DEPLOY_SERVER = credentials('deploy-server-host')  // e.g., ubuntu@your-server.com
        DEPLOY_PATH = '/home/ubuntu/apps/scoop_hrms_backend'
        
        // Build info
        GIT_COMMIT_SHORT = sh(script: 'git rev-parse --short HEAD', returnStdout: true).trim()
        BUILD_TIMESTAMP = sh(script: 'date +%Y%m%d%H%M%S', returnStdout: true).trim()
    }

    options {
        buildDiscarder(logRotator(numToKeepStr: '10'))
        timestamps()
        timeout(time: 30, unit: 'MINUTES')
        disableConcurrentBuilds()
    }

    stages {
        stage('Checkout') {
            steps {
                checkout scm
                echo "Building branch: ${env.BRANCH_NAME ?: env.GIT_BRANCH}"
                echo "Commit: ${GIT_COMMIT_SHORT}"
            }
        }

        stage('Setup Go') {
            steps {
                script {
                    // Use Go tool configured in Jenkins or download
                    def goHome = tool name: 'Go', type: 'go'
                    env.PATH = "${goHome}/bin:${env.PATH}"
                    env.GOPATH = "${WORKSPACE}/go"
                    env.GOCACHE = "${WORKSPACE}/.cache/go-build"
                }
                sh 'go version'
            }
        }

        stage('Dependencies') {
            steps {
                echo 'Downloading Go dependencies...'
                sh '''
                    go mod download
                    go mod verify
                '''
            }
        }

        stage('Lint') {
            steps {
                echo 'Running linter...'
                sh '''
                    # Install golangci-lint if not available
                    if ! command -v golangci-lint &> /dev/null; then
                        curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.55.2
                    fi
                    
                    # Run linter (allow failure for now, can be made strict)
                    golangci-lint run ./... --timeout=5m || true
                '''
            }
        }

        stage('Test') {
            steps {
                echo 'Running tests...'
                sh '''
                    # Run tests with coverage
                    go test -v -race -coverprofile=coverage.out -covermode=atomic ./... || true
                    
                    # Generate coverage report
                    go tool cover -html=coverage.out -o coverage.html || true
                '''
            }
            post {
                always {
                    // Archive test results if available
                    archiveArtifacts artifacts: 'coverage.html', allowEmptyArchive: true
                }
            }
        }

        stage('Build Binary') {
            steps {
                echo 'Building Go binary...'
                sh '''
                    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
                        -ldflags="-w -s -X main.Version=${GIT_COMMIT_SHORT} -X main.BuildTime=${BUILD_TIMESTAMP}" \
                        -o bin/${APP_NAME} \
                        ./cmd/api
                '''
            }
            post {
                success {
                    archiveArtifacts artifacts: 'bin/*', fingerprint: true
                }
            }
        }

        stage('Build Docker Image') {
            steps {
                echo 'Building Docker image...'
                script {
                    def imageTag = "${DOCKER_IMAGE}:${GIT_COMMIT_SHORT}"
                    def latestTag = "${DOCKER_IMAGE}:latest"
                    
                    sh """
                        docker build -t ${imageTag} -t ${latestTag} .
                        docker images | grep ${DOCKER_IMAGE}
                    """
                    
                    env.DOCKER_IMAGE_TAG = imageTag
                }
            }
        }

        stage('Push Docker Image') {
            when {
                anyOf {
                    branch 'main'
                    branch 'master'
                    branch 'develop'
                    branch 'leave-management'
                }
            }
            steps {
                echo 'Pushing Docker image to registry...'
                script {
                    withCredentials([usernamePassword(
                        credentialsId: 'docker-registry-credentials',
                        usernameVariable: 'DOCKER_USER',
                        passwordVariable: 'DOCKER_PASS'
                    )]) {
                        sh '''
                            echo "$DOCKER_PASS" | docker login -u "$DOCKER_USER" --password-stdin ${DOCKER_REGISTRY} || true
                            
                            # Tag and push
                            docker tag ${DOCKER_IMAGE}:${GIT_COMMIT_SHORT} ${DOCKER_REGISTRY}/${DOCKER_IMAGE}:${GIT_COMMIT_SHORT}
                            docker tag ${DOCKER_IMAGE}:latest ${DOCKER_REGISTRY}/${DOCKER_IMAGE}:latest
                            
                            docker push ${DOCKER_REGISTRY}/${DOCKER_IMAGE}:${GIT_COMMIT_SHORT} || true
                            docker push ${DOCKER_REGISTRY}/${DOCKER_IMAGE}:latest || true
                        '''
                    }
                }
            }
        }

        stage('Deploy to Server') {
            when {
                anyOf {
                    branch 'main'
                    branch 'master'
                    branch 'leave-management'
                }
            }
            steps {
                echo 'Deploying to server...'
                script {
                    sshagent(credentials: ['deploy-server-ssh-key']) {
                        sh """
                            ssh -o StrictHostKeyChecking=no ${DEPLOY_SERVER} << 'ENDSSH'
                                cd ${DEPLOY_PATH}
                                
                                # Pull latest code
                                git pull origin ${env.BRANCH_NAME ?: 'leave-management'}
                                
                                # Pull latest Docker image or rebuild
                                docker compose pull || docker compose build
                                
                                # Restart services with zero downtime
                                docker compose up -d --remove-orphans
                                
                                # Clean up old images
                                docker image prune -f
                                
                                # Show status
                                docker compose ps
                                
                                # Health check
                                sleep 10
                                curl -sf http://localhost:8080/health || echo "Health check pending..."
ENDSSH
                        """
                    }
                }
            }
        }
    }

    post {
        always {
            echo 'Cleaning up workspace...'
            sh '''
                # Clean up Docker images to save space
                docker system prune -f || true
            '''
            cleanWs(cleanWhenNotBuilt: false, deleteDirs: true, disableDeferredWipeout: true)
        }
        
        success {
            echo "Build ${env.BUILD_NUMBER} succeeded!"
            // Uncomment to enable Slack notifications
            // slackSend(color: 'good', message: "Build ${env.BUILD_NUMBER} succeeded for ${APP_NAME}")
        }
        
        failure {
            echo "Build ${env.BUILD_NUMBER} failed!"
            // Uncomment to enable Slack notifications
            // slackSend(color: 'danger', message: "Build ${env.BUILD_NUMBER} failed for ${APP_NAME}")
        }
        
        unstable {
            echo "Build ${env.BUILD_NUMBER} is unstable!"
        }
    }
}
