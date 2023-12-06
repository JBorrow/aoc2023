import math

from functools import reduce

PLOT_TEST = False
READ_DATA = True

EPSILON = 1e-10

def load_data(filename: str) -> (list[int], list[int]):
    with open(filename, "r") as file:
        lines = file.readlines()

    times = [int("".join(lines[0].split(":")[1].split()))]
    goals = [int("".join(lines[1].split(":")[1].split()))]

    return times, goals

def bounds(time: int, goal: int) -> (float, float):
    x_low_b = 0.5 * (time - math.sqrt(time * time - 4 * goal))
    x_high_b = 0.5 * (math.sqrt(time * time - 4 * goal) + time)

    return x_low_b, x_high_b

def number_valid(time: int, goal: int) -> int:
    x_low_b, x_high_b = bounds(time, goal)
    return int(math.floor(x_high_b - EPSILON)) - int(math.ceil(x_low_b + EPSILON)) + 1

def number_of_ways_to_win(filename: str) -> int:
    times, goals = load_data(filename)
    return reduce(lambda x, y: x * y, [number_valid(t, g) for t, g in zip(times, goals)])

if READ_DATA:
    import sys

    filename = sys.argv[1]

    print("Number of ways to win: ", number_of_ways_to_win(filename))


if PLOT_TEST:
    if READ_DATA:
        raise Exception("Cannot read data and plot test at the same time")

    import sys

    time = int(sys.argv[1])
    goal = int(sys.argv[2])

    import matplotlib.pyplot as plt

    xs = range(time + 1)
    ys = [(time - x) * x for x in xs]

    x_low_b, x_high_b = bounds(time, goal)
    number_valid = number_valid(time, goal)

    print(x_low_b, x_high_b)

    colors = ["C0" if x_low_b < x < x_high_b else "C1" for x in xs]

    plt.scatter(xs, ys, c=colors)
    plt.axhline(goal, color="k", linestyle="dashed")
    plt.axvline(x_high_b)
    plt.axvline(x_low_b)
    plt.text(0.05, 0.95, f"Number of valid solutions: {number_valid}", transform=plt.gca().transAxes)

    plt.show()
