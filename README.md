# 🔁 Matching Engine

The Matching Engine is responsible for periodically reading offers and ride requests from the database, applying a matching algorithm, and publishing the results to a message queue.

---

## 📚 Matching Algorithm

The matching algorithm is based on the paper:  
**"Optimizing Carpool Scheduling Algorithm through Partition Merging"**  
📄 DOI: [10.1109/ICC.2018.8422976](https://doi.org/10.1109/ICC.2018.8422976)

---

## 📌 Responsibilities & Features

- ⏱ **Scheduled Execution**  
  Runs in **batch mode** at configurable intervals instead of in real-time. This is optimized for recurring or pre-planned trips (e.g., commuting), improving match quality and vehicle occupancy.

- 👥 **Group Requests**  
  Supports multiple riders per request if they share the same origin, destination, and preferences. Enables higher occupancy, helping families or coworkers ride together.

- ➕ **One-to-Many Matching**  
  Matches **multiple requests to a single offer** in a single run, maximizing match rates and reducing unused vehicle capacity.

- ⚙️ **Match Constraints**
  - **Dynamic Capacity**: Riders entering/leaving at different points must never exceed vehicle capacity.
  - **Requests per Driver Limit**: Ensures no more than a configured number of requests per offer are matched.
  - **Preference Compatibility**: Riders and drivers must share compatible trip preferences.
  - **Time Windows**: Pickup must occur after a rider’s earliest departure and drop-off before latest arrival.
  - **Walking Duration**: Includes allowed walking time in timing constraints to improve match chances.
  - **Detour Limits**: Driver's total trip time must not exceed their direct route plus a configured detour margin.

---

## 🔁 Interactions

- 📤 **Publishes results** (successful matches) to a **message queue** for downstream processing
- ⏲ **Triggered periodically** by a **scheduler**, enabling batch-based global optimization

---

## 📣 Maintainers

This service is part of the **3alsekka Carpooling System** (Graduation Project - Alexandria University, 2025).
