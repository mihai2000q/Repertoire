import Order from "../../types/Order.ts";

const artistAlbumsOrders: Order[] = [
  { value: 'release_date desc', label: 'Latest Releases' },
  { value: 'title asc', label: 'Alphabetically' },
  { value: 'release_date asc', label: 'First Releases' },
]

export default artistAlbumsOrders;
