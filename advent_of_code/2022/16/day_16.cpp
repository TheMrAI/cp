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

auto dfs(std::unordered_map<int, Node> const& edgeMap,
  int nodeId,
  int remainingTime,
  int perMinuteRelease,
  // we are cheating here and squeezing the visited into a 64 bit wide unsigned, because we know there are at most only
  // 54 nodes to manage, otherwise this wouldn't work!
  uint64_t seen) -> int
{
  if (remainingTime <= 1 || static_cast<std::size_t>(std::popcount(seen)) == edgeMap.size()) {
    return perMinuteRelease * remainingTime;
  }

  auto released = 0;
  // open current valve if it is worth opening
  if (edgeMap.at(nodeId).flowRate != 0) {
    --remainingTime;
    released = perMinuteRelease;
    perMinuteRelease += edgeMap.at(nodeId).flowRate;
  }
  // mark node as seen
  seen |= (uint64_t{ 1 } << nodeId);

  auto maximumThatCanBeReleased = perMinuteRelease * remainingTime;
  for (auto const& targetNode : edgeMap.at(nodeId).edges) {
    auto targetNodeVisited = (seen & (uint64_t{ 1 } << targetNode.first)) != 0;
    if (targetNodeVisited) { continue; }
    auto travelTime = targetNode.second;
    auto releasedWhileTraversing = travelTime * perMinuteRelease;
    maximumThatCanBeReleased = std::max(maximumThatCanBeReleased,
      releasedWhileTraversing + dfs(edgeMap, targetNode.first, remainingTime - travelTime, perMinuteRelease, seen));
  }
  return released + maximumThatCanBeReleased;
}

// curiously we don't even have to know where the agent actually is, only where they are going and when they arrive
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
  // all node seen or no time left
  if (currentTime == endTime || static_cast<std::size_t>(std::popcount(seen)) == edgeMap.size()) {
    return releasedWhileStepping;
  }

  auto userTargets = std::vector<std::pair<int, int>>{};
  userTargets.emplace_back(user.targetNodeId, user.arrivalTime);
  if (user.arrivalTime == currentTime) {
    auto userNodeSeen = (seen & (uint64_t{ 1 } << user.targetNodeId)) != 0;
    if (!userNodeSeen) {
      seen |= (uint64_t{ 1 } << user.targetNodeId);
      userTargets.emplace_back(user.targetNodeId, currentTime + 1);
      perMinuteRelease += edgeMap.at(user.targetNodeId).flowRate;
    } else {
      // collect candidates
      for (auto const& userCandidate : edgeMap.at(user.targetNodeId).edges) {
        auto candidateSeen = (seen & (uint64_t{ 1 } << userCandidate.first)) != 0;
        if (candidateSeen || userCandidate.first == elephant.targetNodeId) { continue; }
        userTargets.emplace_back(userCandidate.first, currentTime + userCandidate.second);
      }
    }
  }

  auto elephantTargets = std::vector<std::pair<int, int>>{};
  elephantTargets.emplace_back(elephant.targetNodeId, elephant.arrivalTime);
  if (elephant.arrivalTime == currentTime) {
    auto elephantNodeSeen = (seen & (uint64_t{ 1 } << elephant.targetNodeId)) != 0;
    if (!elephantNodeSeen) {
      seen |= (uint64_t{ 1 } << elephant.targetNodeId);
      elephantTargets.emplace_back(elephant.targetNodeId, currentTime + 1);
      perMinuteRelease += edgeMap.at(elephant.targetNodeId).flowRate;
    } else {
      // collect candidates
      for (auto const& elephantCandidate : edgeMap.at(elephant.targetNodeId).edges) {
        auto candidateSeen = (seen & (uint64_t{ 1 } << elephantCandidate.first)) != 0;
        if (candidateSeen || elephantCandidate.first == user.targetNodeId) { continue; }
        elephantTargets.emplace_back(elephantCandidate.first, currentTime + elephantCandidate.second);
      }
    }
  }

  auto additionalMaxReleasablePressure = 0;
  for (auto const& userCandidate : userTargets) {
    for (auto const& elephantCandidate : elephantTargets) {
      auto userTargetArrivalTime =
        currentTime >= userCandidate.second ? endTime : std::min(userCandidate.second, endTime);
      auto elephantTargetArrivalTime =
        currentTime >= elephantCandidate.second ? endTime : std::min(elephantCandidate.second, endTime);
      auto travelTime = std::min(userTargetArrivalTime, elephantTargetArrivalTime) - currentTime;
      additionalMaxReleasablePressure = std::max(additionalMaxReleasablePressure,
        dfs_second_part(edgeMap,
          Agent{ userCandidate.first, userCandidate.second },
          Agent{ elephantCandidate.first, elephantCandidate.second },
          travelTime,
          currentTime + travelTime,
          endTime,
          perMinuteRelease,
          seen));
    }
  }

  return releasedWhileStepping + additionalMaxReleasablePressure;
}

auto main() -> int
{
  ios::sync_with_stdio(false);
  cin.tie(nullptr);

  auto originalGraph = parseGraph();
  auto nodeEncoding = encodeGraphNodes(originalGraph);
  auto originalEdgeMap = toStandardEdgeMap(originalGraph, nodeEncoding);
  auto edgeMap = compressEdgeMap(originalEdgeMap);

  auto totalReleased = dfs(edgeMap, 0, 30, 0, uint64_t{ 0 });
  std::cout << "Total pressure released: " << totalReleased << std::endl;
  // 2595 too low
  auto totalReleasedWithElephantBro =
    dfs_second_part(edgeMap, Agent{ 0, 1 }, Agent{ 0, 1 }, 0, 1, 26, 0, uint64_t{ 1 });
  std::cout << "Total pressure release with elephant companion: " << totalReleasedWithElephantBro << std::endl;

  return 0;
}
