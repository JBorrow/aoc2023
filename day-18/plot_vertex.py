import matplotlib.pyplot as plt
import numpy as np

verticies = np.array([
[6,0],
[6,5],
[4,5],
[4,7],
[6,7],
[6,9],
[1,9],
[1,7],
[0,7],
[0,5],
[2,5],
[2,2],
[0,2],
[0,0],
[6,0],
]).T

plt.plot(verticies[0], verticies[1])
plt.show()
