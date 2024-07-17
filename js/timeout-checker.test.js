// BEGIN: Test for ticker function
test("ticker function should log time left every second", () => {
  // Mock console.log
  const originalLog = console.log;
  console.log = jest.fn();

  // Call the ticker function with a time value of 10 seconds
  ticker(10);

  // Wait for 11 seconds (to allow for 10 iterations of the setInterval)
  jest.advanceTimersByTime(11000);

  // Expect console.log to have been called 10 times
  expect(console.log).toHaveBeenCalledTimes(10);

  // Restore console.log
  console.log = originalLog;
});
// END: Test for ticker function
