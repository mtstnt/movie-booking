import os
from pathlib import Path
import sys

SERVICES = ["gateway", "user", "employee", "movie", "booking", "notification"]

def generate_proto(service_name: str, language: str) -> None:
  for service in SERVICES:
    if os.path.exists(f"./pb/{service}.proto"):
      os.system(f"protoc --{language}_out=./services/{service_name} --{language}_opt=paths=source_relative \
        --{language}-grpc_out=./services/{service_name} --{language}-grpc_opt=paths=source_relative \
        ./pb/{service}.proto")

def gen(target_services: list[str]) -> None:
  for service in target_services:
      if service in SERVICES:
        generate_proto(service, 'go')
      else: print(f"Service {service} is not registered")

def clean(target_services: list[str]) -> None:
  if len(target_services) == 0:
    target_services = SERVICES

  for service in target_services:
    if service in SERVICES:
      p = f"./services/{service}/pb"
      if os.path.isdir(p):
        for entries in os.scandir(p):
          os.unlink(entries.path)
        os.rmdir(p)
      else: print(f"Path {Path(p).absolute} does not exist.")
    else: print(f"Service {service} is not registered.")

def show_help() -> None:
  print("""Proto Generator v0.0.1
  Available commands:
    - gen [services...] -> Generate proto file for services
    - clean [services... \\ null] -> Clean proto/pb directory for services (or all services if not specified)
    - help -> show help page.""")

# Handle generate for specified service
command_args = sys.argv[1:]
if command_args[0].lower() == "gen": 
  gen(command_args[1:])
elif command_args[0].lower() == "clean":
  clean(command_args[1:])
elif command_args[0].lower() == "help":
  show_help()
else:
  print(f"Command {command_args[0]} is not registered. Run with command 'help' to see commands")