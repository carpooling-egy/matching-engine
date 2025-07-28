from fastapi import FastAPI, HTTPException
from pydantic import BaseModel
from typing import List, Tuple
from ortools.constraint_solver import routing_enums_pb2, pywrapcp
import logging

# Set up logging
logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)

app = FastAPI()


class VRPData(BaseModel):
    time_matrix: List[List[int]]
    time_windows: List[Tuple[int, int]]
    no_of_riders_per_request: List[int]
    vehicle_capacity: int
    pickup_and_dropoffs: List[Tuple[int, int]]
    max_route_duration: int
    timeout : int = 100  
    method : str = "parallel_cheapest_insertion"
    enable_guided_local_search: bool = False



@app.post("/solve")
def solve_vrp(data: VRPData):

    logger.info("Received VRP data for solving")

    manager, routing = create_routing_model(data)
    time_callback_idx = routing.RegisterTransitCallback(
        lambda from_idx, to_idx: data.time_matrix[manager.IndexToNode(from_idx)][manager.IndexToNode(to_idx)]
    )
    routing.SetArcCostEvaluatorOfAllVehicles(time_callback_idx)

    # Add time & capacity constraints
    add_constraints(routing, manager, data, time_callback_idx)

    # Configure solver
    search_params = pywrapcp.DefaultRoutingSearchParameters()
    
    if data.method == "path_cheapest_arc":
        search_params.first_solution_strategy = routing_enums_pb2.FirstSolutionStrategy.PATH_CHEAPEST_ARC
    elif data.method == "global_cheapest_arch":
        search_params.first_solution_strategy = routing_enums_pb2.FirstSolutionStrategy.GLOBAL_CHEAPEST_ARC
    elif data.method == "automatic":
        search_params.first_solution_strategy = routing_enums_pb2.FirstSolutionStrategy.AUTOMATIC
    else:
        search_params.first_solution_strategy = routing_enums_pb2.FirstSolutionStrategy.PARALLEL_CHEAPEST_INSERTION
    
    search_params.time_limit.nanos   = data.timeout * 1_000_000 
    search_params.log_search = True

    if data.enable_guided_local_search:
        # Enable guided local search for better optimization
        search_params.local_search_metaheuristic = (
            routing_enums_pb2.LocalSearchMetaheuristic.GUIDED_LOCAL_SEARCH
        )

    solution = routing.SolveWithParameters(search_params)


    if solution:
        route = extract_solution(data, manager, routing, solution)
        return {"success": True, "route": route}
    else:
        logger.warning("No feasible solution found")
        return {"success": False, "route": []}



def create_routing_model(data: VRPData):
    manager = pywrapcp.RoutingIndexManager(
        len(data.time_matrix),  # number of locations
        1,                      # one vehicle
        [0],                    # start depot
        [len(data.time_matrix)-1],                    # end depot
    )
    routing = pywrapcp.RoutingModel(manager)
    return manager, routing


def add_constraints(routing, manager, data: VRPData, time_callback_idx: int):
    routing.AddDimension(
        time_callback_idx,
        0,             # no slack
        data.max_route_duration,          # max route duration
        True,
        "Time"
    )
    time_dim = routing.GetDimensionOrDie("Time")

    # Time windows
    for idx, window in enumerate(data.time_windows):
        if idx in {0, len(data.time_matrix)-1}:  # start/end depot
            continue
        node_index = manager.NodeToIndex(idx)
        time_dim.CumulVar(node_index).SetRange(window[0], window[1])

    

    # Capacity
    capacity_callback_idx = routing.RegisterUnaryTransitCallback(
        lambda idx: data.no_of_riders_per_request[manager.IndexToNode(idx)]
    )
    routing.AddDimensionWithVehicleCapacity(
        capacity_callback_idx,
        0,  # no slack
        [data.vehicle_capacity],
        True,
        "Capacity"
    )

    # Pickups & dropoffs
    for pickup_node, dropoff_node in data.pickup_and_dropoffs:
        pickup_idx = manager.NodeToIndex(pickup_node)
        dropoff_idx = manager.NodeToIndex(dropoff_node)

        routing.AddPickupAndDelivery(pickup_idx, dropoff_idx)
        time_dim = routing.GetDimensionOrDie("Time")

        routing.solver().Add(
            routing.VehicleVar(pickup_idx) == routing.VehicleVar(dropoff_idx)
        )
        routing.solver().Add(time_dim.CumulVar(pickup_idx) <= time_dim.CumulVar(dropoff_idx))



def extract_solution(data: VRPData, manager, routing, solution):
    vehicle_id = 0
    time_dim = routing.GetDimensionOrDie("Time")
    cap_dim = routing.GetDimensionOrDie("Capacity")

    index = routing.Start(vehicle_id)
    route_info = []

    while not routing.IsEnd(index):
        node = manager.IndexToNode(index)
        arrival = solution.Min(time_dim.CumulVar(index))

        route_info.append({
            "node": node,
            "arrival_time": arrival,
        })

        index = solution.Value(routing.NextVar(index))

    # Final stop
    node = manager.IndexToNode(index)
    arrival = solution.Min(time_dim.CumulVar(index))

    route_info.append({
        "node": node,
        "arrival_time": arrival,
    })

    return route_info




def dummy_test():
    
    payload = {
        "time_matrix": [
            [0, 600000000000, 20, 10],
            [10, 0, 10, 20],
            [20, 10, 0, 10],
            [10, 20, 10, 0]
        ],
        "time_windows": [
            [0, 100],  # depot
            [0, 5],  # pickup
            [0, 50],   # dropoff
            [0, 100]   # depot
        ],
        "no_of_riders_per_request": [0, 1, -1, 0],
        "vehicle_capacity": 1
    }

    # try:
    response = solve_vrp(VRPData(**payload))
    print("Response:", response)



if __name__ == "__main__":
    dummy_test()