#include <bits/stdc++.h>

enum class RobotType {
    NONE,
    ORE,
    CLAY,
    OBSIDIAN,
    GEODE
};

// overkill, but it makes me fell better
struct Robot {
    RobotType type;
    uint32_t ore;
    uint32_t clay;
    uint32_t obsidian;
};

struct Blueprint {
    uint32_t id;
    Robot none_robot; // "robot" to easily handle no resource costs
    Robot ore_robot;
    Robot clay_robot;
    Robot obsidian_robot;
    Robot geode_robot;
};

struct ProductionStatus{
    uint32_t ore_total;
    uint32_t clay_total;
    uint32_t obsidian_total;
    uint32_t geode_total;

    uint32_t ore_production;
    uint32_t clay_production;
    uint32_t obsidian_production;
    uint32_t geode_production;
};

auto can_build_robot(Robot const& robot, ProductionStatus const& status) -> bool {
    if (robot.ore <= status.ore_total && robot.clay <= status.clay_total && robot.obsidian <= status.obsidian_total) {
        return true;
    }
    return false;
}

struct SimulationPoint {
    uint32_t minute;
    ProductionStatus status;
};

auto update_status(SimulationPoint simulation_point, Robot const& robot) -> SimulationPoint {
    // start building robot
    simulation_point.status.ore_total -= robot.ore;
    simulation_point.status.clay_total -= robot.clay;
    simulation_point.status.obsidian_total -= robot.obsidian;

    // collect resources
    simulation_point.status.ore_total += simulation_point.status.ore_production;
    simulation_point.status.clay_total += simulation_point.status.clay_production;
    simulation_point.status.obsidian_total += simulation_point.status.obsidian_production;
    simulation_point.status.geode_total += simulation_point.status.geode_production;

    // update robot count
    switch (robot.type) {
        case RobotType::NONE: {
            break;
        }
        case RobotType::ORE: {
            simulation_point.status.ore_production += 1;
            break;
        }
        case RobotType::CLAY: {
            simulation_point.status.clay_production += 1;
            break;
        }
        case RobotType::OBSIDIAN: {
            simulation_point.status.obsidian_production += 1;
            break;
        }
        case RobotType::GEODE: {
            simulation_point.status.geode_production += 1;
            break;
        }
    }
    ++simulation_point.minute;
    return simulation_point;
}

auto sum_with_arithmetic_progression(uint32_t start, uint32_t remaining_steps) -> uint32_t {
    auto n = start + remaining_steps;
    return (start + n) * n / 2;
}

// search as BFT and prune as many branches as we can
auto get_maximal_geode_count(Blueprint const& blueprint, uint32_t const simulation_minutes, SimulationPoint simulate_from = SimulationPoint{1, ProductionStatus{0, 0, 0, 0, 1, 0, 0, 0}}, uint32_t const seen_maximum = 0) -> uint32_t {
    uint32_t maximum_necessary_ore_production = std::max({blueprint.ore_robot.ore, blueprint.clay_robot.ore, blueprint.obsidian_robot.ore, blueprint.geode_robot.ore});
    uint32_t maximum_necessary_clay_production = std::max({blueprint.ore_robot.clay, blueprint.clay_robot.clay, blueprint.obsidian_robot.clay, blueprint.geode_robot.clay});
    uint32_t maximum_necessary_obsidian_production = std::max({blueprint.ore_robot.obsidian, blueprint.clay_robot.obsidian, blueprint.obsidian_robot.obsidian, blueprint.geode_robot.obsidian});
    
    uint32_t maximal_geode_count = seen_maximum;

    std::queue<SimulationPoint> options;
    options.push(std::move(simulate_from));

    while(!options.empty()) {
        auto to_evaluate = options.front();
        options.pop();

        auto remaining_minutes = simulation_minutes - to_evaluate.minute;
        
        auto absolute_maximal_output = to_evaluate.status.geode_total + sum_with_arithmetic_progression(to_evaluate.status.geode_production, remaining_minutes);
        // if we assumed we could produce a geode robot every round from here on, would we beat the maximum, if no
        // this branch is useless
        if (to_evaluate.minute > simulation_minutes || absolute_maximal_output < maximal_geode_count) {
            continue;
        }
        // if we could produce one geode bot each round, we can calculate the maximum achievable globally
        // , because we are searching using BFT, so the first layer we find such a match is the best one
        // TODO fix this, should give nice speedup
        // if (to_evaluate.status.ore_production >= blueprint.geode_robot.ore && to_evaluate.status.clay_production >= blueprint.geode_robot.clay && to_evaluate.status.obsidian_production >= blueprint.geode_robot.obsidian) {
        //     maximal_geode_count = absolute_maximal_output;
        //     // std::cout << "Quick exit at minute: " << to_evaluate.minute << std::endl;
        //     break;
        // }

        // build no robot, this is always true
        if (can_build_robot(blueprint.none_robot, to_evaluate.status)) {
            auto after_update = update_status(to_evaluate, blueprint.none_robot);
            maximal_geode_count = std::max(maximal_geode_count, after_update.status.geode_total);
            options.push(after_update);
        }
        // don't build more production robot for a resource that cannot be consumed in a round
        if (can_build_robot(blueprint.ore_robot, to_evaluate.status) && to_evaluate.status.ore_production < maximum_necessary_ore_production) {
            auto after_update = update_status(to_evaluate, blueprint.ore_robot);
            maximal_geode_count = std::max(maximal_geode_count, after_update.status.geode_total);
            options.push(after_update);
        } 
        if (can_build_robot(blueprint.clay_robot, to_evaluate.status) && to_evaluate.status.clay_production < maximum_necessary_clay_production) {
            auto after_update = update_status(to_evaluate, blueprint.clay_robot);
            maximal_geode_count = std::max(maximal_geode_count, after_update.status.geode_total);
            options.push(after_update);
        } 
        if (can_build_robot(blueprint.obsidian_robot, to_evaluate.status) && to_evaluate.status.obsidian_production < maximum_necessary_obsidian_production) {
            auto after_update = update_status(to_evaluate, blueprint.obsidian_robot);
            // maximal_geode_count = std::max(maximal_geode_count, after_update.status.geode_total);
            // options.push(after_update);
            maximal_geode_count = std::max({maximal_geode_count, after_update.status.geode_total, get_maximal_geode_count(blueprint, simulation_minutes, after_update, maximal_geode_count)});
        }
        if (can_build_robot(blueprint.geode_robot, to_evaluate.status)) {
            auto after_update = update_status(to_evaluate, blueprint.geode_robot);
            // maximal_geode_count = std::max(maximal_geode_count, after_update.status.geode_total);
            // options.push(after_update);
            maximal_geode_count = std::max({maximal_geode_count, after_update.status.geode_total, get_maximal_geode_count(blueprint, simulation_minutes, after_update, maximal_geode_count)});
        }
    }

    return maximal_geode_count;
}

auto main() -> int {
    std::ios::sync_with_stdio(false);
    std::cin.tie(nullptr);

    auto line = std::string{};
    auto data_regex = std::regex(R"(Blueprint (\d+): Each ore robot costs (\d+) ore. Each clay robot costs (\d+) ore. Each obsidian robot costs (\d+) ore and (\d+) clay. Each geode robot costs (\d+) ore and (\d+) obsidian.)");

    auto total_score = uint32_t{0};

    while (std::getline(std::cin, line)) {
        auto match = std::smatch{};
        if (!std::regex_match(line, match, data_regex)) {
            std::cerr << "Regular expression failure on input: " << line << std::endl;
        }

        auto blueprint_id = static_cast<uint32_t>(std::stoi(match[1]));
        auto none_robot = Robot{RobotType::NONE, 0, 0, 0};
        auto ore_robot = Robot{RobotType::ORE, static_cast<uint32_t>(std::stoi(match[2])), 0, 0};
        auto clay_robot = Robot{RobotType::CLAY, static_cast<uint32_t>(std::stoi(match[3])), 0, 0};
        auto obsidian_robot = Robot{RobotType::OBSIDIAN, static_cast<uint32_t>(std::stoi(match[4])), static_cast<uint32_t>(std::stoi(match[5])), 0};
        auto geode_robot = Robot{RobotType::GEODE, static_cast<uint32_t>(std::stoi(match[6])), 0, static_cast<uint32_t>(std::stoi(match[7]))};

        auto blueprint = Blueprint{
            blueprint_id,
            none_robot,
            ore_robot,
            clay_robot,
            obsidian_robot,
            geode_robot
        };

        auto geode_count = get_maximal_geode_count(blueprint, 32);
        auto blueprint_score = blueprint_id * geode_count;
        total_score += blueprint_score;
        std::cout << blueprint_id << " geode count: " << geode_count << " ,with score: " << blueprint_score << std::endl;
    }
    std::cout << "Total score: " << total_score << std::endl;

    return 0;
}


// First run on 32 min input that didn't explode with out of memory (Still used up more than 10 GiBs though)!!!

// mrai@potatoland:~/tinker/cp/advent_of_code/2022$ time ./day_19 < day_19_2.in
// 1 geode count: 30 ,with score: 30
// 2 geode count: 4 ,with score: 8
// 3 geode count: 17 ,with score: 51
// Total score: 89

// real    3m6.314s
// user    3m1.291s
// sys     0m4.990s