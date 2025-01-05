terraform {
    required_providers {
        kubernetes = {
            source = "hashicorp/kubernetes"
            version = "2.30.0"
        }
        helm = {
            source = "hashicorp/helm"
            version = "2.13.2"
        }
    }
}

provider "helm" {
    kubernetes {
        config_path = "~/.kube/config"
        config_context = "minikube"
    }
}

provider "kubernetes" {
    config_path = "~/.kube/config"
    config_context = "minikube"
}

resource "kubernetes_namespace" "spirrel" {
    metadata {
        name = "spirrel"
    }
}

resource "kubernetes_namespace" "elastic-system" {
    metadata {
        name = "elastic-system"
    }
}


resource "helm_release" "eck" {
    name       = "eck"
    chart      = "eck-operator"
    repository = "https://helm.elastic.co"
    namespace  = "elastic-system"
    version    = "2.16.0"
    depends_on = [ 
        kubernetes_namespace.elastic-system
    ]
}

resource "kubernetes_manifest" "elasticsearch" {
    manifest = {
        "apiVersion" = "elasticsearch.k8s.elastic.co/v1"
        "kind" = "Elasticsearch"
        "metadata" = {
            "name" = "elasticsearch"
            "namespace" = "spirrel"
        }
        "spec" = {
            "version" = "8.17.0"
            "nodeSets" = [
                {
                    "name" = "default",
                    "count" = 1,
                    "config" = {
                        "node.store.allow_mmap" = "false"
                    }
                }
            ]
        }
    }
    depends_on = [ 
        helm_release.eck
    ]
}

resource "kubernetes_manifest" "kibana" {
    manifest = {
        "apiVersion" = "kibana.k8s.elastic.co/v1"
        "kind" = "Kibana"
        "metadata" = {
            "name" = "kibana"
            "namespace" = "spirrel"
        }
        "spec" = {
            "version" = "8.17.0"
            "count" = 1
            "elasticsearchRef" = {
                "name" = "elasticsearch"
            }
        }
    }
    depends_on = [ 
        helm_release.eck
    ]
}

resource "helm_release" "example" {
    name       = "spirrel"
    chart      = "${path.module}/../charts"
    namespace  = kubernetes_namespace.spirrel.metadata.0.name
    values = [
        <<-EOT
            elasticsearch:
                host: https://elasticsearch-es-http:9200
                apiKey: WmF1Rk41UUJEeVdlNlpVQ1NxaXg6NXdkeDljc0xTdXVTdjlWSWhJSV9pdw==
            image:
                repository: ghcr.io/chazapp/spirrel
                tag: 1.0.1
                
        EOT
    ]
}