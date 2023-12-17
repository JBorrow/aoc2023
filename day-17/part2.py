"""
Try again in python
"""

import sys
from heapq import heappop, heappush

input = sys.stdin.read().strip().split("\n")

start = (0, 0)
end = (len(input[0]) - 1, len(input) - 1)


def run(min_dist, max_dist):
    # (x, y, cost, direction)
    q = [
        (int(0), *start, (1, 0)),
        (int(0), *start, (0, -1)),
        (int(0), *start, (0, 1)),
        (int(0), *start, (0, -1)),
    ]

    seen = set()
    costs = {}

    while q:
        cost, x, y, direction = heappop(q)

        if (x, y) == end:
            return cost

        if (x, y, direction) in seen:
            continue

        seen.add((x, y, direction))

        for new_direction in [(0, 1), (1, 0), (0, -1), (-1, 0)]:
            if (new_direction[0] == -direction[0] and new_direction[1] == -direction[1]):
                # Not allowed
                continue

            if (new_direction == direction):
                # Not allowed
                continue

            new_cost = 0
            new_x = x
            new_y = y

            # Add on one for each step
            for distance in range(1, max_dist + 1):
                new_x += new_direction[0]
                new_y += new_direction[1]

                if new_x < 0 or new_y < 0 or new_x >= len(input[0]) or new_y >= len(input):
                    break

                new_cost += int(input[new_y][new_x])

                if distance < min_dist:
                    continue

                if costs.get((new_x, new_y, new_direction), float("inf")) <= cost + new_cost:
                    continue

                costs[(new_x, new_y, new_direction)] = cost + new_cost

                heappush(q, (cost + new_cost, new_x, new_y, new_direction))

    return -1



print(run(4, 10))