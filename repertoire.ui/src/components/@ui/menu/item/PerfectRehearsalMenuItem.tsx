import { IconChecks } from '@tabler/icons-react'
import { toast } from 'react-toastify'
import { useAddPerfectSongRehearsalMutation } from '../../../../state/api/songsApi.ts'
import MenuItemConfirmation from './MenuItemConfirmation.tsx'
import { useAddPerfectRehearsalsToArtistsMutation } from '../../../../state/api/artistsApi.ts'
import { useAddPerfectRehearsalsToAlbumsMutation } from '../../../../state/api/albumsApi.ts'
import { useAddPerfectRehearsalsToPlaylistsMutation } from '../../../../state/api/playlistsApi.ts'

interface PerfectRehearsalMenuItemProps {
  id: string
  closeMenu: () => void
  type: 'artist' | 'album' | 'song' | 'playlist'
}

function PerfectRehearsalMenuItem({ id, closeMenu, type }: PerfectRehearsalMenuItemProps) {
  const [addPerfectRehearsalsToArtists, { isLoading: isArtistLoading }] =
    useAddPerfectRehearsalsToArtistsMutation()
  const [addPerfectRehearsalsToAlbums, { isLoading: isAlbumLoading }] =
    useAddPerfectRehearsalsToAlbumsMutation()
  const [addPerfectSongRehearsal, { isLoading: isSongLoading }] =
    useAddPerfectSongRehearsalMutation()
  const [addPerfectRehearsalsToPlaylists, { isLoading: isPlaylistLoading }] =
    useAddPerfectRehearsalsToPlaylistsMutation()
  const isLoading = isArtistLoading || isAlbumLoading || isSongLoading || isPlaylistLoading

  async function handleAddPerfectRehearsal() {
    switch (type) {
      case 'artist':
        await addPerfectRehearsalsToArtists({ ids: [id] }).unwrap()
        break
      case 'album':
        await addPerfectRehearsalsToAlbums({ ids: [id] }).unwrap()
        break
      case 'song':
        await addPerfectSongRehearsal({ id: id }).unwrap()
        break
      case 'playlist':
        await addPerfectRehearsalsToPlaylists({ ids: [id] }).unwrap()
        break
    }
    toast.success(`Perfect rehearsal added!`)
    closeMenu()
  }

  return (
    <MenuItemConfirmation
      isLoading={isLoading}
      onConfirm={handleAddPerfectRehearsal}
      leftSection={<IconChecks size={14} />}
    >
      Perfect Rehearsal
    </MenuItemConfirmation>
  )
}

export default PerfectRehearsalMenuItem
