import Order from '../../types/Order.ts'
import ArtistProperty from '../../utils/enums/ArtistProperty.ts'
import OrderType from '../../utils/enums/OrderType.ts'

const artistsOrders: Order[] = [
  { label: 'Name', property: ArtistProperty.Name, checked: false },
  { label: 'Albums', property: ArtistProperty.Albums, checked: false },
  { label: 'Songs', property: ArtistProperty.Songs, checked: false },
  { label: 'Band Members', property: ArtistProperty.BandMembers, checked: false },
  { label: 'Rehearsals', property: ArtistProperty.Rehearsals, checked: false },
  { label: 'Confidence', property: ArtistProperty.Confidence, checked: false },
  { label: 'Progress', property: ArtistProperty.Progress, checked: false },
  {
    label: 'Creation Date',
    property: ArtistProperty.CreationDate,
    type: OrderType.Descending,
    checked: true
  },
  { label: 'Last Modified', property: ArtistProperty.LastModified, checked: false }
]

export default artistsOrders
