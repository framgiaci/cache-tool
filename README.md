# How to install

```bash
sudo curl -o /usr/local/bin/cache-tool https://raw.githubusercontent.com/framgiaci/cache-tool/master/dist/cache-tool && \ 
sudo chmod +x /usr/local/bin/cache-tool

// Or local build (Make sure you installed golang on your machine)

go build -o cache-tool main.go
mv cache-tool /usr/local/bin/
sudo chmod +x /usr/local/bin/cache-tool

```

# How to use

```
// To create cache - at current folder, make sure you have framgia-ci.yml and run:
cache-tool --create

// To restore cache
cache-tool --restore
```
