
pipeline {
    agent { label 'Debian-jk-agent' }

    options {
        timestamps()
        disableConcurrentBuilds()
    }

    environment {
        APP_NAME = 'project-devops'
        DEPLOYMENT = 'project-devops'
        CONTAINER = 'project-devops'
        NAMESPACE = "${env.BRANCH_NAME == 'main' ? 'prod' : 'dev'}"
    }

    stages {
        stage('Verificar entorno') {
            steps {
                sh '''
                    whoami
                    docker --version
                    minikube status
                    kubectl get nodes
                '''
            }
        }

        stage('Construir imagen Docker') {
            steps {
                sh '''
                    docker build -t ${APP_NAME}:${BUILD_NUMBER} .
                '''
            }
        }

        stage('Cargar imagen en Minikube') {
            steps {
                sh '''
                    minikube image load ${APP_NAME}:${BUILD_NUMBER}
                '''
            }
        }

        stage('Desplegar en Kubernetes') {
            steps {
                sh '''
                    sed "s|image: project-devops:1.0|image: ${APP_NAME}:${BUILD_NUMBER}|" \
                        k8s/deployment.yaml > deployment-rendered.yaml

                    kubectl apply -f deployment-rendered.yaml
                    kubectl apply -f k8s/service.yaml

                    kubectl rollout status deployment/${DEPLOYMENT} --timeout=180s
                '''
            }
        }

        stage('Verificar réplicas') {
            steps {
                sh '''
                    kubectl wait --for=condition=available \
                        deployment/${DEPLOYMENT} --timeout=180s

                    kubectl get deployments
                    kubectl get pods -o wide
                    kubectl get services

                    READY=$(kubectl get deployment ${DEPLOYMENT} \
                        -o jsonpath='{.status.readyReplicas}')

                    test "$READY" = "2"
                '''
            }
        }
    }

    post {
        success {
            echo 'Despliegue completado correctamente.'
        }
        failure {
            echo 'Error en el pipeline. Revisar la consola.'
        }
    }
}
