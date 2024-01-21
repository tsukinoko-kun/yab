require("env")

os.execute("go build ./...")
os.execute("go test ./... -v")
