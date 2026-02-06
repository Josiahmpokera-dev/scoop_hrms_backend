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
                    env.GIT_BRANCH_NAME = sh(script: 'git rev-parse --abbrev-ref HEAD || echo "leave-management"', returnStdout: true).trim()
                }
                echo "Building branch: ${env.GIT_BRANCH_NAME}"
                echo "Commit: ${env.GIT_COMMIT_SHORT}"
            }
        }

        stage('Build Docker Image') {
            steps {
                echo 'Building Docker image...'
                sh '''
                    echo "Building image: ${DOCKER_IMAGE}:${GIT_COMMIT_SHORT}"
                    docker build \
                        -t ${DOCKER_IMAGE}:${GIT_COMMIT_SHORT} \
                        -t ${DOCKER_IMAGE}:latest \
                        .
                    echo "=== Built Images ==="
                    docker images | grep ${DOCKER_IMAGE} | head -5
                '''
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
            }
            steps {
                script {
                    def hasCredentials = false
                    try {
                        withCredentials([usernamePassword(credentialsId: 'docker-registry-credentials', usernameVariable: 'U', passwordVariable: 'P')]) {
                            hasCredentials = true
                        }
                    } catch (Exception e) {
                        echo "Docker registry credentials not configured - skipping push"
                    }
                    
                    if (hasCredentials) {
                        withCredentials([
                            string(credentialsId: 'docker-registry-url', variable: 'REGISTRY_URL'),
                            usernamePassword(credentialsId: 'docker-registry-credentials', usernameVariable: 'DOCKER_USER', passwordVariable: 'DOCKER_PASS')
                        ]) {
                            sh '''
                                echo "${DOCKER_PASS}" | docker login -u "${DOCKER_USER}" --password-stdin ${REGISTRY_URL}
                                docker tag ${DOCKER_IMAGE}:latest ${REGISTRY_URL}/${DOCKER_IMAGE}:${GIT_COMMIT_SHORT}
                                docker tag ${DOCKER_IMAGE}:latest ${REGISTRY_URL}/${DOCKER_IMAGE}:latest
                                docker push ${REGISTRY_URL}/${DOCKER_IMAGE}:${GIT_COMMIT_SHORT}
                                docker push ${REGISTRY_URL}/${DOCKER_IMAGE}:latest
                                docker logout ${REGISTRY_URL}
                            '''
                        }
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
                script {
                    def canDeploy = false
                    try {
                        withCredentials([sshUserPrivateKey(credentialsId: 'deploy-server-ssh-key', keyFileVariable: 'K')]) {
                            canDeploy = true
                        }
                    } catch (Exception e) {
                        echo "Deploy SSH credentials not configured - skipping deployment"
                    }
                    
                    if (canDeploy) {
                        withCredentials([
                            string(credentialsId: 'deploy-server-host', variable: 'DEPLOY_SERVER'),
                            string(credentialsId: 'deploy-server-path', variable: 'DEPLOY_PATH')
                        ]) {
                            sshagent(credentials: ['deploy-server-ssh-key']) {
                                sh '''
                                    ssh -o StrictHostKeyChecking=no ${DEPLOY_SERVER} "
                                        cd ${DEPLOY_PATH} &&
                                        git pull origin ${GIT_BRANCH_NAME} &&
                                        docker compose build &&
                                        docker compose up -d --remove-orphans &&
                                        docker image prune -f &&
                                        docker compose ps &&
                                        sleep 5 &&
                                        curl -sf http://localhost:8080/health || echo 'Health check pending...'
                                    "
                                '''
                            }
                        }
                    }
                }
            }
        }
    }

    post {
        always {
            echo "Build #${env.BUILD_NUMBER} completed with status: ${currentBuild.currentResult}"
        }
        success {
            echo 'Build successful!'
            sh 'docker system prune -f || true'
        }
        failure {
            echo 'Build failed!'
        }
        cleanup {
            cleanWs()
        }
    }
}
