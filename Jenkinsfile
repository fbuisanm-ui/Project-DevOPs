
pipeline {
    agent { label 'Debian-jk-agent' }

    options {
        timestamps()
        disableConcurrentBuilds()
        skipDefaultCheckout(false)
    }

    environment {
        APP_NAME = 'project-devops'
        DEPLOYMENT = 'project-devops'
        CONTAINER = 'project-devops'
        NAMESPACE = "${env.BRANCH_NAME == 'main' ? 'prod' : 'dev'}"
        IMAGE_TAG = "${env.BRANCH_NAME}-${env.BUILD_NUMBER}"
    }

    stages {
        stage('Validar rama y entorno') {
            steps {
                script {
                    if (!(env.BRANCH_NAME in ['develop', 'main'])) {
                        error("Rama no autorizada: ${env.BRANCH_NAME}")
                    }
                }

                sh '''
                    echo "Rama: ${BRANCH_NAME}"
                    echo "Namespace: ${NAMESPACE}"
                    echo "Imagen: ${APP_NAME}:${IMAGE_TAG}"

                    whoami
                    git --version
                    docker --version
                    minikube status
                    kubectl get nodes
                '''
            }
        }

        stage('Verificar namespace') {
            steps {
                sh '''
                    kubectl create namespace "${NAMESPACE}" \
                        --dry-run=client -o yaml | kubectl apply -f -
                '''
            }
        }

        stage('Construir imagen Docker') {
            steps {
                sh '''
                    docker build \
                        -t "${APP_NAME}:${IMAGE_TAG}" .
                '''
            }
        }

        stage('Cargar imagen en Minikube') {
            steps {
                sh '''
                    minikube image load "${APP_NAME}:${IMAGE_TAG}"
                '''
            }
        }

        stage('Desplegar en Kubernetes') {
            steps {
                sh '''
                    sed \
                        "s|image: project-devops:1.0|image: ${APP_NAME}:${IMAGE_TAG}|" \
                        k8s/deployment.yaml > deployment-rendered.yaml

                    grep -q \
                        "image: ${APP_NAME}:${IMAGE_TAG}" \
                        deployment-rendered.yaml

                    kubectl apply -n "${NAMESPACE}" \
                        -f deployment-rendered.yaml

                    kubectl apply -n "${NAMESPACE}" \
                        -f k8s/service.yaml

                    kubectl rollout status \
                        deployment/"${DEPLOYMENT}" \
                        -n "${NAMESPACE}" \
                        --timeout=180s
                '''
            }
        }

        stage('Verificar despliegue') {
            steps {
                sh '''
                    kubectl wait \
                        --for=condition=available \
                        deployment/"${DEPLOYMENT}" \
                        -n "${NAMESPACE}" \
                        --timeout=180s

                    kubectl get deployments -n "${NAMESPACE}"
                    kubectl get pods -n "${NAMESPACE}" -o wide
                    kubectl get services -n "${NAMESPACE}"

                    READY=$(kubectl get deployment "${DEPLOYMENT}" \
                        -n "${NAMESPACE}" \
                        -o jsonpath='{.status.readyReplicas}')

                    test "${READY}" = "2"

                    DEPLOYED_IMAGE=$(kubectl get deployment "${DEPLOYMENT}" \
                        -n "${NAMESPACE}" \
                        -o jsonpath='{.spec.template.spec.containers[0].image}')

                    test "${DEPLOYED_IMAGE}" = "${APP_NAME}:${IMAGE_TAG}"
                '''
            }
        }
    }

    post {
        success {
            echo "Despliegue completado: ${env.BRANCH_NAME} → ${env.NAMESPACE}"
        }

        failure {
            echo 'Error en el pipeline. Revisar la consola.'
        }
    }
}
