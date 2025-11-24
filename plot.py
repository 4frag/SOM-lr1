#!.venv/bin/python

import sys
import pandas as pd
import matplotlib.pyplot as plt

orbit_file = sys.argv[1]
radius_file = sys.argv[2]
speed_file = sys.argv[3]

# --- ORBIT X-Y ---
orbit = pd.read_csv(orbit_file)
plt.figure(figsize=(6,6))
plt.plot(orbit["x"], orbit["y"])
plt.title("Orbit projection (X-Y)")
plt.xlabel("x (km)")
plt.ylabel("y (km)")
plt.axis("equal")
plt.grid()
plt.savefig("orbit_xy.png", dpi=200)
plt.close()

# --- RADIUS ---
radius = pd.read_csv(radius_file)
plt.figure()
plt.plot(radius["t"] / 3600, radius["r"])
plt.title("Radius vs Time")
plt.xlabel("Time (hours)")
plt.ylabel("Radius (km)")
plt.grid()
plt.savefig("radius_plot.png", dpi=200)
plt.close()

# --- SPEED ---
speed = pd.read_csv(speed_file)
plt.figure()
plt.plot(speed["t"] / 3600, speed["speed"])
plt.title("Speed vs Time")
plt.xlabel("Time (hours)")
plt.ylabel("Speed (km/s)")
plt.grid()
plt.savefig("speed_plot.png", dpi=200)
plt.close()

# --- 3D ORBIT ---
fig = plt.figure()
ax = fig.add_subplot(111, projection='3d')
ax.plot(orbit["x"], orbit["y"], orbit["z"])
ax.set_title("3D Orbit")
ax.set_xlabel("x")
ax.set_ylabel("y")
ax.set_zlabel("z")
plt.savefig("orbit_3d.png", dpi=200)
plt.close()

print("Plots saved: orbit_xy.png, radius_plot.png, speed_plot.png, orbit_3d.png")
