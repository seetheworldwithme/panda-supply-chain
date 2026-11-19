import unittest
from contract_compare import decode, differences


class CompareTests(unittest.TestCase):
    def test_numeric_string_is_not_number(self):
        self.assertTrue(differences(decode('"1"'), decode('1')))

    def test_boolean_is_not_integer(self):
        self.assertTrue(differences(decode('true'), decode('1')))

    def test_decimal_scale_is_preserved(self):
        self.assertTrue(differences(decode('1.00'), decode('1')))

    def test_null_is_not_absent(self):
        self.assertTrue(differences(decode('{"x":null}'), decode('{}')))

    def test_array_order_is_preserved(self):
        self.assertTrue(differences(decode('[1,2]'), decode('[2,1]')))

    def test_object_key_order_does_not_matter(self):
        self.assertFalse(differences(decode('{"x":1,"y":null}'), decode('{"y":null,"x":1}')))


if __name__ == '__main__':
    unittest.main()
