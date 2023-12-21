import sys
import attr

input = sys.stdin.read().strip().split("\n")

DEBUG = False

all_nodes = {}


def debug_print(x):
    if DEBUG:
        print(x)


PULSES = {True: 0, False: 0}
QUEUE = []

PART2 = True

@attr.s
class Conjunction:
    name: str = attr.ib()
    outputs: list[str] = attr.ib()
    inputs: dict[bool] = attr.ib()

    def __call__(self, name: str, high: bool):
        PULSES[high] += 1
        debug_print(f"&calling {self.name} from {name} with {high}")

        if name not in self.inputs:
            raise AttributeError

        self.inputs[name] = high

        return_value = not all(self.inputs.values())

        for output in self.outputs:
            QUEUE.append((output, self.name, return_value))


@attr.s
class FlipFlop:
    name: str = attr.ib()
    outputs: list[str] = attr.ib()
    state: bool = attr.ib(default=False)

    def __call__(self, name: str, high: bool):
        PULSES[high] += 1
        debug_print(f"%calling {self.name} from {name} with {high}")
        if high:
            return
        else:
            self.state = not self.state

        for output in self.outputs:
            QUEUE.append((output, self.name, self.state))


@attr.s
class Broadcaster:
    name: str = attr.ib()
    outputs: list[str] = attr.ib()

    def __call__(self, name: str, high: bool):
        PULSES[high] += 1
        debug_print(f"bcalling {self.name} from {name} with {high}")
        for output in self.outputs:
            QUEUE.append((output, self.name, high))

@attr.s
class Output:
    name: str = attr.ib()

    def __call__(self, name: str, high: bool):
        PULSES[high] += 1
        debug_print(f"ocalling {self.name} from {name} with {high}")

        if PART2 and self.name == "rx":
            print("rx", high)
            if not high:
                exit()
            return


get_name = lambda x: x.split("->")[0].strip().replace("%", "").replace("&", "")
contains_output = lambda x, line: x in line.split("->")[1]

for line in input:
    name = get_name(line)
    outputs = [x.strip() for x in line.split("->")[1].split(",")]

    if "broadcaster" in line:
        all_nodes[name] = Broadcaster(name, outputs)
    elif "&" in line:
        all_nodes[name] = Conjunction(
            name,
            outputs,
            {get_name(x): False for x in input if contains_output(name, x)},
        )
    elif "%" in line:
        all_nodes[name] = FlipFlop(name, outputs)


for line in input:
    outputs = [x.strip() for x in line.split("->")[1].split(",")]

    for output in outputs:
        if output not in all_nodes:
            all_nodes[output] = Output(output)
    

if not PART2:
    for i in range(1000):
        debug_print(f"push {i}")
        all_nodes["broadcaster"]("start", False)

        while QUEUE:
            name, caller, high = QUEUE.pop(0)
            all_nodes[name](caller, high)

        debug_print(all_nodes)
        debug_print("")

    debug_print(PULSES)
    print("total", PULSES[True] * PULSES[False])

if PART2:
    i = 0
    while True:
        i += 1
        print(f"push {i}")

        all_nodes["broadcaster"]("start", False)

        while QUEUE:
            name, caller, high = QUEUE.pop(0)
            all_nodes[name](caller, high)
