#include <bits/stdc++.h>
#include <format>

using namespace std;
using i6 = int64_t;

template<typename Output, typename Iter>
  requires std::derived_from<Output, std::ostream> && std::weakly_incrementable<Iter> && std::indirectly_readable<Iter>
auto dump_to(Output& output, Iter begin, Iter end)
{
  while (begin != end) {
    output << *begin;
    if (std::next(begin) != end) { output << ' '; }
    ++begin;
  }
}

// Source: https://codeforces.com/blog/entry/78852
template<typename T>
  requires std::integral<T>
auto ceil2(T a, T b) -> T
{
  if (a == 0) return 0;
  return ((a - 1) / b) + 1;
}

struct Valve
{
  int flowRate;
  std::vector<std::string> tunnels;
};

auto parseGraph() -> std::unordered_map<std::string, Valve>
{
  auto graph = std::unordered_map<std::string, Valve>{};

  auto line = string{};
  auto data_regex = regex(R"(Valve (\w{2}) has flow rate=(\d+); tunnels? leads? to valves? (.+))");
  auto targets_regex = regex(R"([A-Z]{2})");

  while (getline(cin, line)) {
    auto match = std::smatch{};
    if (std::regex_match(line, match, data_regex)) {
      std::cout << "Match: " << match[1].str() << " flow rate: " << match[2].str() << " path to: " << match[3].str()
                << std::endl;

      auto valve = Valve{};

      valve.flowRate = std::atoi(match[2].str().c_str());
      auto targets = match[3].str();
      auto targets_begin = std::sregex_iterator{ targets.cbegin(), targets.cend(), targets_regex };
      auto targets_end = std::sregex_iterator{};

      for (auto iter = targets_begin; iter != targets_end; ++iter) { valve.tunnels.push_back(iter->str()); }
      graph[match[1].str()] = valve;
    } else {
      std::cerr << "Regex match failure" << std::endl;
    }
  }
  return graph;
}

auto encodeGraphNodes(std::unordered_map<std::string, Valve> const& inputGraph) -> std::unordered_map<std::string, int>
{
  auto nodeIds = std::vector<std::string>(inputGraph.size(), "");
  auto i = 0;
  for (auto const& entry : inputGraph) {
    nodeIds[i] = entry.first;
    ++i;
  }
  ranges::sort(nodeIds);
  auto encoding = std::unordered_map<std::string, int>{};
  for (auto i = std::size_t{ 0 }; i < nodeIds.size(); ++i) { encoding[nodeIds[i]] = i; }
  return encoding;
}

struct Node
{
  int flowRate;
  std::unordered_map<int, int> edges;// nodeId => weight
};

auto toStandardEdgeMap(std::unordered_map<std::string, Valve> const inputGraph,
  std::unordered_map<std::string, int> const encoding) -> std::unordered_map<int, Node>
{
  auto edgeMap = std::unordered_map<int, Node>{};
  for (auto const& entry : inputGraph) {
    auto nodeId = encoding.at(entry.first);
    auto node = Node{};
    node.flowRate = entry.second.flowRate;
    for (auto const& neighbor : entry.second.tunnels) { node.edges.insert({ encoding.at(neighbor), 1 }); }
    edgeMap[nodeId] = node;
  }

  return edgeMap;
}

auto floydWarshall(std::unordered_map<int, Node> const edgeMap) -> std::vector<std::vector<int>>
{
  constexpr auto INF = std::numeric_limits<int>::max();

  auto distanceMatrix = std::vector<std::vector<int>>(edgeMap.size(), std::vector<int>(edgeMap.size(), INF));

  for (auto const& sourceNode : edgeMap) {
    distanceMatrix[sourceNode.first][sourceNode.first] = 0;
    for (auto const& targetNode : sourceNode.second.edges) {
      distanceMatrix[sourceNode.first][targetNode.first] = targetNode.second;
    }
  }

  for (auto k = std::size_t{ 0 }; k < edgeMap.size(); ++k) {
    for (auto i = std::size_t{ 0 }; i < edgeMap.size(); ++i) {
      for (auto j = std::size_t{ 0 }; j < edgeMap.size(); ++j) {
        if (distanceMatrix[i][k] == INF || distanceMatrix[k][j] == INF) { continue; }
        distanceMatrix[i][j] = std::min(distanceMatrix[i][j], distanceMatrix[i][k] + distanceMatrix[k][j]);
      }
    }
  }

  return distanceMatrix;
}

auto compressEdgeMap(std::unordered_map<int, Node> const& edgeMap) -> std::unordered_map<int, Node>
{
  // Floyd-Warshall
  auto shortestPaths = floydWarshall(edgeMap);

  auto compressedEdgeMap = std::unordered_map<int, Node>{};
  // add all positive flow rate nodes
  for (auto const& node : edgeMap) {
    if (node.second.flowRate > 0) {
      compressedEdgeMap[node.first] = node.second;
      // clear all connections, they will get reconnected later
      compressedEdgeMap[node.first].edges = std::unordered_map<int, int>{};
    }
  }
  // connect all positive flow nodes to all other
  for (auto& node : compressedEdgeMap) {
    for (auto const& targetNode : compressedEdgeMap) {
      if (node.first == targetNode.first) { continue; }
      node.second.edges[targetNode.first] = shortestPaths[node.first][targetNode.first];
    }
  }
  // add zero id node, because that is the start even if it has zero flow
  compressedEdgeMap[0] = edgeMap.at(0);
  compressedEdgeMap[0].edges = std::unordered_map<int, int>{};
  // connect to all other nodes
  for (auto const& targetNode : compressedEdgeMap) {
    if (targetNode.first == 0) { continue; }
    compressedEdgeMap[0].edges[targetNode.first] = shortestPaths[0][targetNode.first];
  }

  return compressedEdgeMap;
}

// find all optimal paths for one agent
// The trick we are gonna employ here is that we merge, the travel and valve release
// times into one range of time. This way we don't have to worry about how to handle
// the cases where one agent has arrived somewhere but hasn't opened the valve yet.
auto dfs(std::unordered_map<int, Node> const& edgeMap,
  int nodeId,
  int arrivalTime,
  int simulationStep,
  int currentTime,
  int endTime,
  int perMinuteRelease,
  int releasedSoFar,
  // we are cheating here and squeezing the visited into a 64 bit wide unsigned, because we know there are at most only
  // 54 nodes to manage, otherwise this wouldn't work!
  uint64_t seen,
  std::unordered_map<uint64_t, int>& bestPaths) -> void
{
  auto releasedWhileStepping = perMinuteRelease * simulationStep;
  releasedSoFar += releasedWhileStepping;
  auto remainingSteps = endTime - currentTime;
  // always open the valves the agent has arrived at
  // at the start node the flow rate may be zero and we skip that
  if (arrivalTime == currentTime && edgeMap.at(nodeId).flowRate != 0) {
    seen |= (uint64_t{ 1 } << nodeId);
    perMinuteRelease += edgeMap.at(nodeId).flowRate;
  }
  // if didn't move any more, this would be the additionally max released pressure
  auto additionalReleasedPressure = perMinuteRelease * remainingSteps;
  // insert/update best release for seen path, if we didn't make any more moves
  bestPaths[seen] = std::max(bestPaths[seen], releasedSoFar + additionalReleasedPressure);

  // all node seen or no time left
  if (currentTime == endTime || static_cast<std::size_t>(std::popcount(seen)) == edgeMap.size()) { return; }

  // pick a node and go boiii
  for (auto const& entry : edgeMap.at(nodeId).edges) {
    auto nodeSeen = (seen & (uint64_t{ 1 } << entry.first)) != 0;
    if (nodeSeen) { continue; }
    auto stepCount = entry.second + 1;// travel time + valve opening time
    simulationStep = std::min({ stepCount, remainingSteps });
    dfs(edgeMap,
      entry.first,
      stepCount + currentTime,
      simulationStep,
      simulationStep + currentTime,
      endTime,
      perMinuteRelease,
      releasedSoFar,
      seen,
      bestPaths);
  }
}

// Part two
// After reading a few ideas, I have found one that I have liked.
// Here it is how it goes: To find an optimal path we have to essentially DFS the whole graph and collect what is the
// best score we can get. This means that we inevitably have to check all paths either way. If we identify a path
// through the graph using the bitmask of the visited nodes, we can collect the maximum releasable pressure for each
// such path, into a map. Additionally after every visited and valve opened node we have to store the case where the
// agent decided not to search any more. This way we generate all the possible combinations and for each we can
// store the best possible route for a single agent.
// Then, we observe that for a disjoint path pair, we can simply add their best scores together and find the optimal
// result for two agents. This is due to the nature of the data we have collected. If two paths are disjoint, and they
// are optimal for a given period of time, we should just be able to add them together.
// Find the biggest such pair and we are done.

// Total running time for part one and two together:
// real    0m0.105s
// user    0m0.095s
// sys     0m0.010s
// Not gonna lie, I am impressed. Given how complex the problem is, it seems solving it may not be that expensive.
// Initially I didn't even want to consider these path collection approaches as there would be nodeCount! paths to
// check. However, it seems that most of those paths get pruned away and even on the big input we would only have a
// total
// of only 9903 paths to check, which is a very small amount.
auto main() -> int
{
  ios::sync_with_stdio(false);
  cin.tie(nullptr);

  auto originalGraph = parseGraph();
  auto nodeEncoding = encodeGraphNodes(originalGraph);
  auto originalEdgeMap = toStandardEdgeMap(originalGraph, nodeEncoding);
  auto edgeMap = compressEdgeMap(originalEdgeMap);

  std::cout << "Total number of non-zero valves: " << edgeMap.size() - 1 << std::endl;

  // do not mark the first node as visited, or part two won't be able to find disjoint paths
  auto bestReleaseForPaths = std::unordered_map<uint64_t, int>{};
  dfs(edgeMap, 0, 0, 0, 0, 30, 0, 0, uint64_t{ 0 }, bestReleaseForPaths);
  std::cout << "Total number of optimal paths: " << bestReleaseForPaths.size() << std::endl;
  auto optimalReleaseForOneAgent = bestReleaseForPaths.begin()->second;
  for (auto const& entry : bestReleaseForPaths) {
    optimalReleaseForOneAgent = std::max(optimalReleaseForOneAgent, entry.second);
  }
  std::cout << "Total pressure released: " << optimalReleaseForOneAgent << std::endl;
  // 2595 too low
  // 2715 not good
  auto bestReleaseForPathsTwoAgents = std::unordered_map<uint64_t, int>{};
  dfs(edgeMap, 0, 0, 0, 0, 26, 0, 0, uint64_t{ 0 }, bestReleaseForPathsTwoAgents);
  auto bestFlowForTwoAgents = 0;
  for (auto const& agentOne : bestReleaseForPathsTwoAgents) {
    for (auto const& agentTwo : bestReleaseForPathsTwoAgents) {
      if ((agentOne.first & agentTwo.first) != 0) { continue; }
      bestFlowForTwoAgents = std::max(bestFlowForTwoAgents, agentOne.second + agentTwo.second);
    }
  }
  // auto totalReleasedWithElephantBro =
  //   dfs_second_part(edgeMap, Agent{ 0, 0 }, Agent{ 0, 0 }, 0, 0, 30, 0, uint64_t{ 1 });
  std::cout << "Total pressure release with elephant companion: " << bestFlowForTwoAgents << std::endl;

  return 0;
}


// My 3rd attempt to implement simultaneous simulation. Doesn't work, though...
// Something is missing that I just can't see.
// The trick we are gonna employ here is that we merge, the travel and valve release
// times into one range of time. This way we don't have to worry about how to handle
// the cases where one agent has arrived somewhere but hasn't opened the valve yet.
struct Agent
{
  int targetNodeId;
  int arrivalTime;
};

auto dfs_second_part(std::unordered_map<int, Node> const& edgeMap,
  Agent user,
  Agent elephant,
  int simulationStep,
  int currentTime,
  int endTime,
  int perMinuteRelease,
  uint64_t seen) -> int
{
  auto releasedWhileStepping = perMinuteRelease * simulationStep;
  auto remainingSteps = endTime - currentTime;

  // always open the valves the agent has arrived at
  // , if it isn't already open (the opening time has been calculated into the travel time)
  if (user.arrivalTime == currentTime && !((seen & (uint64_t{ 1 } << user.targetNodeId)) != 0)) {
    seen |= (uint64_t{ 1 } << user.targetNodeId);
    perMinuteRelease += edgeMap.at(user.targetNodeId).flowRate;
  }
  if (elephant.arrivalTime == currentTime && !((seen & (uint64_t{ 1 } << elephant.targetNodeId)) != 0)) {
    seen |= (uint64_t{ 1 } << elephant.targetNodeId);
    perMinuteRelease += edgeMap.at(elephant.targetNodeId).flowRate;
  }
  // if nobody moved this round, this would be the additionally max releasable pressure
  auto additionalMaxReleasablePressure = perMinuteRelease * remainingSteps;

  // all node seen or no time left
  if (currentTime == endTime || static_cast<std::size_t>(std::popcount(seen)) == edgeMap.size()) {
    return releasedWhileStepping + additionalMaxReleasablePressure;
  }

  // only user steps, if it can
  if (user.arrivalTime <= currentTime) {
    for (auto const& forUser : edgeMap.at(user.targetNodeId).edges) {
      auto userSeen = (seen & (uint64_t{ 1 } << forUser.first)) != 0;
      if (userSeen) { continue; }
      auto userStep = forUser.second + 1;// travel time + valve opening time
      auto elephantStep = std::max(elephant.arrivalTime - currentTime, remainingSteps);
      simulationStep = std::min({ userStep, elephantStep, remainingSteps });
      additionalMaxReleasablePressure = std::max(additionalMaxReleasablePressure,
        dfs_second_part(edgeMap,
          Agent{ forUser.first, userStep + currentTime },
          elephant,
          simulationStep,
          currentTime + simulationStep,
          endTime,
          perMinuteRelease,
          seen));
    }
  }
  // only elephant steps, if it can
  if (elephant.arrivalTime <= currentTime) {
    for (auto const& forElephant : edgeMap.at(elephant.targetNodeId).edges) {
      auto elephantSeen = (seen & (uint64_t{ 1 } << forElephant.first)) != 0;
      if (elephantSeen) { continue; }
      auto elephantStep = forElephant.second + 1;// travel time + opening time
      auto userStep = std::max(user.arrivalTime - currentTime, remainingSteps);
      simulationStep = std::min({ userStep, elephantStep, remainingSteps });
      additionalMaxReleasablePressure = std::max(additionalMaxReleasablePressure,
        dfs_second_part(edgeMap,
          user,
          Agent{ forElephant.first, elephantStep + currentTime },
          simulationStep,
          currentTime + simulationStep,
          endTime,
          perMinuteRelease,
          seen));
    }
  }
  // both try to step
  if (user.arrivalTime <= currentTime && elephant.arrivalTime <= currentTime
      && (edgeMap.size() - static_cast<std::size_t>(std::popcount(seen)) >= 2)) {
    // choose targets
    for (auto const& forUser : edgeMap.at(user.targetNodeId).edges) {
      auto userSeen = (seen & (uint64_t{ 1 } << forUser.first)) != 0;
      if (userSeen) { continue; }
      auto userStep = forUser.second + 1;
      for (auto const& forElephant : edgeMap.at(elephant.targetNodeId).edges) {
        auto elephantSeen = (seen & (uint64_t{ 1 } << forElephant.first)) != 0;
        if (forUser.first == forElephant.first || elephantSeen) { continue; }
        auto elephantStep = forElephant.second + 1;
        simulationStep = std::min({ userStep, elephantStep, remainingSteps });
        additionalMaxReleasablePressure = std::max(additionalMaxReleasablePressure,
          dfs_second_part(edgeMap,
            Agent{ forUser.first, userStep + currentTime },
            Agent{ forElephant.first, elephantStep + currentTime },
            simulationStep,
            currentTime + simulationStep,
            endTime,
            perMinuteRelease,
            seen));
      }
    }
  }

  return releasedWhileStepping + additionalMaxReleasablePressure;
}