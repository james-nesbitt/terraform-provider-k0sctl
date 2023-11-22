# install Mirantis products using parametrized k0sctl
resource "k0sctl_config" "example" {
  metadata {
    name = "test"
  }
  spec {
    k0s {
    }

    host {
      role = "manager"
      ssh {
        address  = "manager1.example.org"
        key_path = "./key.pem"
        user     = "ubuntu"
      }

      hooks {
        apply {
          before = ["ls -la", "pwd"]
        }
      }
      mcr_config {
        debug = true
        bip   = "172.20.0.1/16"

        default_address_pool {
          base = "172.20.0.0/16"
          size = 16
        }
      }
    }

    host {
      role = "worker"
      ssh {
        address  = "worker1.example.org"
        key_path = "./key.pem"
        user     = "ubuntu"
      }
    }

    host {
      role = "worker"
      winrm {
        address  = "windowsworker1.example.org"
        user     = "ubuntu"
        password = "my-win-password"
      }
    }

    host {
      role = "msr"
      ssh {
        address  = "msr1.example.org"
        key_path = "./key.pem"
        user     = "ubuntu"
      }
    }
  }
}
