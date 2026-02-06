pipeline {
    agent any

    environment {
        APP_NAME = 'hrms-backend'
        DOCKER_IMAGE = 'hrms-backend'
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
                script {
                    env.GIT_COMMIT_SHORT = sh(script: 'git rev-parse --short HEAD', returnStdout: true).trim()
                    env.GIT_BRANCH_NAME = sh(script: 'git rev-parse --abbrev-ref HEAD', returnStdout: true).trim()
                }
                echo "Building branch: ${env.GIT_BRANCH_NAME}"
                echo "Commit: ${env.GIT_COMMIT_SHORT}"
            }
        }

        stage('Build Docker Image') {
            steps {
                echo 'Building Docker image...'
                sh """
                    docker build -t ${DOCKER_IMAGE}:${env.GIT_COMMIT_SHORT} -t ${DOCKER_IMAGE}:latest .
                    docker images | grep ${DOCKER_IMAGE}
                """
            }
        }

        stage('Push to Registry') {
            when {
                anyOf {
                    branch 'main'
                    branch 'master'
                    branch 'develop'
                    branch 'leave-management'
                }
                // Only run if credentials exist
                expression {
                    try {
                        withCredentials([usernamePassword(credentialsId: 'docker-registry-credentials', usernameVariable: 'U', passwordVariable: 'P')]) {
                            return true
                        }
                    } catch (Exception e) {
                        echo "Docker registry credentials not configured - skipping push"
                        return false
                    }
                }
            }
            steps {
                echo 'Pushing Docker image to registry...'
                withCredentials([
                    string(credentialsId: 'docker-registry-url', variable: 'REGISTRY_URL'),
                    usernamePassword(credentialsId: 'docker-registry-credentials', usernameVariable: 'DOCKER_USER', passwordVariable: 'DOCKER_PASS')
                ]) {
                    sh '''
                        echo "$DOCKER_PASS" | docker login -u "$DOCKER_USER" --password-stdin ${REGISTRY_URL}
                        docker tag ${DOCKER_IMAGE}:latest ${REGISTRY_URL}/${DOCKER_IMAGE}:${GIT_COMMIT_SHORT}
                        docker tag ${DOCKER_IMAGE}:latest ${REGISTRY_URL}/${DOCKER_IMAGE}:latest
                        docker push ${REGISTRY_URL}/${DOCKER_IMAGE}:${GIT_COMMIT_SHORT}
                        docker push ${REGISTRY_URL}/${DOCKER_IMAGE}:latest
                    '''
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
                // Only run if SSH credentials exist
                expression {
                    try {
                        withCredentials([sshUserPrivateKey(credentialsId: 'deploy-server-ssh-key', keyFileVariable: 'K')]) {
                            return true
                        }
                    } catch (Exception e) {
                        echo "Deploy SSH credentials not configured - skipping deployment"
                        return false
                    }
                }
            }
            steps {
                echo 'Deploying to server...'
                withCredentials([
                    string(credentialsId: 'deploy-server-host', variable: 'DEPLOY_SERVER'),
                    string(credentialsId: 'deploy-server-path', variable: 'DEPLOY_PATH')
                ]) {
                    sshagent(credentials: ['deploy-server-ssh-key']) {
                        sh """
                            ssh -o StrictHostKeyChecking=no ${DEPLOY_SERVER} << 'ENDSSH'
                                cd ${DEPLOY_PATH}
                                git pull origin ${env.GIT_BRANCH_NAME}
                                docker compose build
                                docker compose up -d --remove-orphans
                                docker image prune -f
                                docker compose ps
                                sleep 5
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
            echo 'Pipeline completed'
        }
        success {
            echo "Build ${env.BUILD_NUMBER} succeeded!"
            // Clean up Docker images on success
            sh 'docker system prune -f || true'
        }
        failure {
            echo "Build ${env.BUILD_NUMBER} failed!"
        }
        cleanup {
            cleanWs()
        }
    }
}
