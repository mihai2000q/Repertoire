export default function plural(param: unknown[] | number | undefined) {
  if (!param) return ''

  return typeof param === 'number' ? (param === 1 ? '' : 's') : param.length === 1 ? '' : 's'
}
