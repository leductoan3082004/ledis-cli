# ledis-cli

---

### Installation
- Ensure that `Go` has installed on your machine, if not please refer to this [official document](https://go.dev/doc/install) to install it.
- Run the server by run the `make build` command in your terminal, if `make` has not installed yet u can install it. The produced binary will be called `ledis-cli` located at your root folder.
- If you want to use it as `ledis-cli --host 127.0.0.1 --port 6379` instead of `./ledis-cli --host 127.0.0.1 --port 6379`, you can export the binary to your `PATH` as follow: 

`export PATH="$PATH:/home/youruser/ledisclifolder/ledis-cli"`

### Guide

- After make sure your binary is ready. Please connect to the `Ledis server` by using this command: `ledis-cli --host 127.0.0.1 --port 6379`, replace the host and port in correspondingly if you run `Ledis server` from remote host than your local. If you run your `Ledis server` on local, just run `ledis-cli` is enough, default host is your local and port is 6379.
- Then if connect successfully, you will be prompted into `ledis-cli terminal` and at here you can type your command as guided in `Ledis server` documentation.