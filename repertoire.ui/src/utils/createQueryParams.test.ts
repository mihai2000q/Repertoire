import createQueryParams from "./createQueryParams.ts";

describe('create Query Params', () => {
  it('should create query params', () => {
    const input = {
      id: '1',
      name: 'John',
      searchBy: ['description IS NOT NULL', 'sku IS NOT NULL'],
      orderBy: ['name ASC', 'created_at DESC'],
    }

    let expected = `?id=${input.id}` +
      `&name=${input.name}` +
      `&searchBy=${input.searchBy[0]}` +
      `&searchBy=${input.searchBy[1]}` +
      `&orderBy=${input.orderBy[0]}` +
      `&orderBy=${input.orderBy[1]}`
    expected = expected.replaceAll(' ', '%20')

    const result = createQueryParams(input)

    expect(result).toBe(expected)
  })

  it('should return empty array when the object was empty', () => {
    const result = createQueryParams({})

    expect(result).toBe('')
  })
})
