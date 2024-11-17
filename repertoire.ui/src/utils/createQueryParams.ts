export default function createQueryParams(obj: object): string {
  const queryString = Object.entries(obj)
    .flatMap(([key, value]) =>
      Array.isArray(value)
        ? value.map((val) => `${encodeURIComponent(key)}=${encodeURIComponent(val)}`)
        : value !== undefined && value !== null
          ? `${encodeURIComponent(key)}=${encodeURIComponent(value)}`
          : []
    )
    .join('&')

  return queryString ? `?${queryString}` : ''
}
