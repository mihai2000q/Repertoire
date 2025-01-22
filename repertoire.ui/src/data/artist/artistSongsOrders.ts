import Order from '../../types/Order.ts'

const artistSongsOrders: Order[] = [
  { value: 'release_date desc, title asc', label: 'Latest Releases' },
  { value: 'title asc', label: 'Alphabetically' },
  { value: 'release_date asc, title asc', label: 'First Releases' }
]

export default artistSongsOrders
