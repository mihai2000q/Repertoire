import createFormData from "./createFormData.ts";

describe('create Form Data', () => {
  it('should create form data', () => {
    const obj = {
      id: '1',
      name: 'John',
    }

    const result = createFormData(obj)
    
    expect(result.get('id')).toBe(obj.id)
    expect(result.get('name')).toBe(obj.name)
  })
})
