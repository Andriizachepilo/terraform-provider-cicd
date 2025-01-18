const { add } = require('./app');
const { expect } = require('@jest/globals');

test('adds 3 + 5 to equal 8', () => {
  expect(add(3, 5)).toBe(8);
});
