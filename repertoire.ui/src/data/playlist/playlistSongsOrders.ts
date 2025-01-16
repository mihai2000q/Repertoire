import Order from '../../types/Order.ts'

const playlistSongsOrders: Order[] = [
  { value: 'song_track_no asc', label: 'Track Number' },
  { value: 'release_date desc', label: 'Latest Releases' },
  { value: 'title asc', label: 'Alphabetically' },
  { value: 'release_date asc', label: 'First Releases' }
]

export default playlistSongsOrders
