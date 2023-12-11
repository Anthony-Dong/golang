#include "times.h"
#include "utils.h"
#include <iostream>
#include <spdlog/spdlog.h>

int main() {
    std::cout << test::times::version() << "\n";
    std::cout << test::utils::version() << "\n";
    spdlog::info("hello world");
}