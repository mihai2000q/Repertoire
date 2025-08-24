import Order from '../../types/Order.ts'
import ArtistProperty from '../../types/enums/properties/ArtistProperty.ts'
import OrderType from '../../types/enums/OrderType.ts'
import HomeTopEntity from '../../types/enums/HomeTopEntity.ts'
import AlbumProperty from '../../types/enums/properties/AlbumProperty.ts'
import SongProperty from '../../types/enums/properties/SongProperty.ts'

const homeTopOrderEntities: Map<HomeTopEntity, Order[]> = new Map<HomeTopEntity, Order[]>([
  [
    HomeTopEntity.Artists,
    [
      {
        label: 'Progress',
        property: ArtistProperty.Progress,
        type: OrderType.Descending,
        thenBy: [{ property: ArtistProperty.Name }]
      },
      {
        label: 'Confidence',
        property: ArtistProperty.Confidence,
        type: OrderType.Descending,
        thenBy: [{ property: ArtistProperty.Name }]
      },
      {
        label: 'Rehearsals',
        property: ArtistProperty.Rehearsals,
        type: OrderType.Descending,
        thenBy: [{ property: ArtistProperty.Name }]
      }
    ]
  ],

  [
    HomeTopEntity.Albums,
    [
      {
        label: 'Progress',
        property: AlbumProperty.Progress,
        type: OrderType.Descending,
        thenBy: [{ property: AlbumProperty.Title }]
      },
      {
        label: 'Confidence',
        property: AlbumProperty.Confidence,
        type: OrderType.Descending,
        thenBy: [{ property: AlbumProperty.Title }]
      },
      {
        label: 'Rehearsals',
        property: ArtistProperty.Rehearsals,
        type: OrderType.Descending,
        thenBy: [{ property: AlbumProperty.Title }]
      }
    ]
  ],

  [
    HomeTopEntity.Songs,
    [
      {
        label: 'Progress',
        property: SongProperty.Progress,
        type: OrderType.Descending,
        thenBy: [{ property: SongProperty.Title }]
      },
      {
        label: 'Confidence',
        property: SongProperty.Confidence,
        type: OrderType.Descending,
        thenBy: [{ property: SongProperty.Title }]
      },
      {
        label: 'Rehearsals',
        property: SongProperty.Rehearsals,
        type: OrderType.Descending,
        thenBy: [{ property: SongProperty.Title }]
      }
    ]
  ]
])

export const defaultHomeTopOrderEntities = new Map<HomeTopEntity, Order>([
  [
    HomeTopEntity.Artists,
    {
      label: 'Progress',
      property: ArtistProperty.Progress,
      type: OrderType.Descending,
      thenBy: [{ property: ArtistProperty.Name }]
    }
  ],

  [
    HomeTopEntity.Albums,
    {
      label: 'Progress',
      property: AlbumProperty.Progress,
      type: OrderType.Descending,
      thenBy: [{ property: AlbumProperty.Title }]
    }
  ],

  [
    HomeTopEntity.Songs,
    {
      label: 'Progress',
      property: SongProperty.Progress,
      type: OrderType.Descending,
      thenBy: [{ property: SongProperty.Title }]
    }
  ]
])

export default homeTopOrderEntities
