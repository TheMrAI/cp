#include <bits/stdc++.h>

using namespace std;
using i6 = int64_t;

// Source: https://codeforces.com/blog/entry/78852
int ceil2(i6 a, i6 b)
{
  if (a == 0) return 0;
  return ((a - 1) / b) + 1;
}

struct Monkey
{
  i6 id{ 0 };
  queue<i6> items{};
  char operation{ '+' };
  optional<i6> modify_by;
  i6 divisible_by{ 0 };
  i6 on_true_to{ 0 };
  i6 on_false_to{ 0 };
};

auto main() -> int
{
  ios::sync_with_stdio(false);
  cin.tie(nullptr);

  auto line = string{};
  auto number_regexp = regex(R"((\d+))");
  auto monkeys = vector<Monkey>{};
  do {
    auto monkey_specs = vector<string>(6);
    auto monkey = Monkey{};
    for (auto i = 0; i < 6; ++i) {
      getline(cin, line);
      // cout << line << endl;
      switch (i) {
      case 0: {
        sscanf(line.c_str(), "Monkey %ld:", &monkey.id);
        break;
      }
      case 1: {
        auto items_begin = sregex_iterator(line.begin(), line.end(), number_regexp);
        auto items_end = sregex_iterator();
        for (auto entry = items_begin; entry != items_end; ++entry) { monkey.items.push(stoi(entry->str())); }
        break;
      }
      case 2: {
        auto modify_by = 0;
        auto assigned = sscanf(line.c_str(), " Operation: new = old %c %d", &monkey.operation, &modify_by);
        if (assigned == 2) { monkey.modify_by.emplace(modify_by); }
        break;
      }
      case 3: {
        sscanf(line.c_str(), " Test: divisible by %ld", &monkey.divisible_by);
        break;
      }
      case 4: {
        sscanf(line.c_str(), " If true: throw to monkey %ld:", &monkey.on_true_to);
        break;
      }
      case 5: {
        sscanf(line.c_str(), " If false: throw to monkey %ld", &monkey.on_false_to);
        break;
      }
      }
    }
    monkeys.push_back(move(monkey));
    // cout << "Monkey " << monkey.id << endl;
    // cout << "Items "<< endl;
    // while(!monkey.items.empty()) {
    //     cout << monkey.items.front() << " ";
    //     monkey.items.pop();
    // }
    // cout<<endl;
    // cout << "Op " << monkey.operation << " " << monkey.modify_by.value_or(-1) << endl;
    // cout << "Divisible by " << monkey.divisible_by << endl;
    // cout << "On true " << monkey.on_true_to << endl;
    // cout << "On false " << monkey.on_false_to << endl;
  } while (getline(cin, line));

  auto inspect_counter = vector<uint64_t>(monkeys.size(), 0);
  auto modulus = 1;
  // Chinese remainder theorem
  // We set the modulus such that taking the modulo of worry levels
  // would produce the same remainder for all divisible_by values.
  // Here we can simply multiply all divisible_by values together because
  // all are primes, i.e. they are coprimes for all pairs.
  for (auto const& monkey : monkeys) { modulus *= monkey.divisible_by; }

  for (auto i = 0; i < 10'000; ++i) {
    for (auto mi = 0; mi < monkeys.size(); ++mi) {
      while (!monkeys[mi].items.empty()) {
        ++inspect_counter[mi];

        auto worry_lvl = monkeys[mi].items.front();
        monkeys[mi].items.pop();
        auto modify_by = worry_lvl;
        if (monkeys[mi].modify_by.has_value()) { modify_by = monkeys[mi].modify_by.value(); }

        if (monkeys[mi].operation == '+') {
          // worry_lvl = (worry_lvl + modify_by) / 3;
          worry_lvl = (worry_lvl + modify_by) % modulus;
        } else if (monkeys[mi].operation == '*') {
          // worry_lvl = (worry_lvl * modify_by) / 3;
          worry_lvl = (worry_lvl * modify_by) % modulus;
        }

        if (worry_lvl % monkeys[mi].divisible_by == 0) {
          monkeys[monkeys[mi].on_true_to].items.push(worry_lvl);
        } else {
          monkeys[monkeys[mi].on_false_to].items.push(worry_lvl);
        }
      }
    }
  }
  sort(inspect_counter.begin(), inspect_counter.end());
  cout << inspect_counter[inspect_counter.size() - 1] << endl;
  cout << inspect_counter[inspect_counter.size() - 2] << endl;
  cout << "Monkey business: "
       << inspect_counter[inspect_counter.size() - 1] * inspect_counter[inspect_counter.size() - 2] << endl;

  return 0;
}
