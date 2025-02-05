# test_math_operations.py
import unittest
from math_operations import add

class TestMathOperations(unittest.TestCase):
    def test_add(self):
        # Test that the add function works correctly
        self.assertEqual(add(1, 2), 3)  # 1 + 2 should be 3
        self.assertEqual(add(-1, 1), 0)  # -1 + 1 should be 0
        self.assertEqual(add(0, 0), 0)   # 0 + 0 should be 0

if __name__ == '__main__':
    unittest.main()
