
pipeline {
    agent {
        label 'minikube-local'
    }

    options {
        timestamps()
        disableConcurrentBuilds()
    }

    environment {
        APP_NAME = 'project-devops'
        DEPLOYMENT = 'project-devops'
        CONTAINER = 'project-devops'
    }

    stages {

        stage('Verificar entorno') {
            steps {
                sh '''
                    echo "=== Verificando herramientas ==="
                    docker --version
                    minikube version
                    kubectl version --client

                    echo "=== Estado de Minikube ==="
                    minikube status
                    kubectl get nodes
                '''
            }
        }

        stage('Build Docker') {
            steps {
                sh '''
                    echo "=== Construyendo imagen ==="
                    docker build \
                      -t ${APP_NAME}:${BUILD_NUMBER} .
                '''
            }
        }

        stage('Cargar imagen en Minikube') {
            steps {
                sh '''
                    echo "=== Cargando imagen ==="
                    minikube image load \
                      ${APP_NAME}:${BUILD_NUMBER}
                '''
            }
        }

        stage('Desplegar Kubernetes') {
            steps {
                sh '''
                    echo "=== Aplicando manifiestos ==="
                    kubectl apply -f k8s/

                    echo "=== Actualizando imagen ==="
                    kubectl set image \
                      deployment/${DEPLOYMENT} \
                      ${CONTAINER}=${APP_NAME}:${BUILD_NUMBER}

                    echo "=== Esperando despliegue ==="
                    kubectl rollout status \
                      deployment/${DEPLOYMENT} \
                      --timeout=180s
                '''
            }
        }

        stage('Verificar réplicas') {
            steps {
                sh '''
                    echo "=== Comprobando réplicas ==="

                    kubectl wait \
                      --for=condition=available \
                      deployment/${DEPLOYMENT} \
                      --timeout=180s

                    kubectl get deployments
                    kubectl get pods -o wide
                    kubectl get services

                    echo "=== Verificando 2 réplicas ==="

                    READY=$(kubectl get deployment \
                      ${DEPLOYMENT} \
                      -o jsonpath='{.status.readyReplicas}')

                    if [ "$READY" != "2" ]; then
                        echo "ERROR: No hay 2 réplicas listas"
                        exit 1
                    fi

                    echo "OK: 2 réplicas disponibles"
                '''
            }
        }
    }

    post {
        success {
            echo 'Pipeline completado correctamente.'
        }

        failure {
            echo 'ERROR: Revisar logs del pipeline.'
        }

        always {
            echo 'Finalizada la ejecución del pipeline.'
        }
    }
}
