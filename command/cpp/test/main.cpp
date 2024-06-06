#include "times.h"
#include "utils.h"
#include <iostream>

int main() {
    std::cout << test::times::version() << "\n";
    std::cout << test::utils::version() << "\n";
}