import plural from './plural.ts'

describe('plural', () => {
  it.each([0, 2, 3])("should return 's' when the number is not 1", (input) => {
    const result = plural(input)

    expect(result).toBe('s')
  })

  it('should return empty string when the number is 1', () => {
    const result = plural(1)

    expect(result).toBe('')
  })

  it.each([[[]], [[0, 1]], [[0, 1, 2]]])(
    "should return 's' when the length of the array is not 1",
    (input) => {
      const result = plural(input)

      expect(result).toBe('s')
    }
  )

  it('should return empty string when the length of the array is 1', () => {
    const result = plural([0])

    expect(result).toBe('')
  })

  it('should return empty string when the input is undefined', () => {
    const result = plural(undefined)

    expect(result).toBe('')
  })
})
