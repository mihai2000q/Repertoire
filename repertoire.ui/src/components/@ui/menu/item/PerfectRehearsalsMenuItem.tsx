import { IconChecks } from '@tabler/icons-react'
import { toast } from 'react-toastify'
import { useAddPerfectSongRehearsalsMutation } from '../../../../state/api/songsApi.ts'
import MenuItemConfirmation from './MenuItemConfirmation.tsx'
import { useAddPerfectRehearsalsToArtistsMutation } from '../../../../state/api/artistsApi.ts'
import { useAddPerfectRehearsalsToAlbumsMutation } from '../../../../state/api/albumsApi.ts'
import { useAddPerfectRehearsalsToPlaylistsMutation } from '../../../../state/api/playlistsApi.ts'

interface PerfectRehearsalsMenuItemProps {
  ids: string[]
  onClose: () => void
  type: 'artists' | 'albums' | 'songs' | 'playlists'
}

function PerfectRehearsalsMenuItem({ ids, onClose, type }: PerfectRehearsalsMenuItemProps) {
  const [addPerfectRehearsalsToArtists, { isLoading: isArtistsLoading }] =
    useAddPerfectRehearsalsToArtistsMutation()
  const [addPerfectRehearsalsToAlbums, { isLoading: isAlbumsLoading }] =
    useAddPerfectRehearsalsToAlbumsMutation()
  const [addPerfectSongRehearsals, { isLoading: isSongsLoading }] =
    useAddPerfectSongRehearsalsMutation()
  const [addPerfectRehearsalsToPlaylists, { isLoading: isPlaylistsLoading }] =
    useAddPerfectRehearsalsToPlaylistsMutation()
  const isLoading = isArtistsLoading || isAlbumsLoading || isSongsLoading || isPlaylistsLoading

  async function handleAddPerfectRehearsals() {
    switch (type) {
      case 'artists':
        await addPerfectRehearsalsToArtists({ ids: ids }).unwrap()
        break
      case 'albums':
        await addPerfectRehearsalsToAlbums({ ids: ids }).unwrap()
        break
      case 'songs':
        await addPerfectSongRehearsals({ ids: ids }).unwrap()
        break
      case 'playlists':
        await addPerfectRehearsalsToPlaylists({ ids: ids }).unwrap()
        break
    }
    toast.success(`Perfect rehearsals added!`)
    onClose()
  }

  return (
    <MenuItemConfirmation
      isLoading={isLoading}
      onConfirm={handleAddPerfectRehearsals}
      leftSection={<IconChecks size={14} />}
    >
      Perfect Rehearsals
    </MenuItemConfirmation>
  )
}

export default PerfectRehearsalsMenuItem
