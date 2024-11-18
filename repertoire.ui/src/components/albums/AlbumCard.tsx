import Album from '../../types/models/Album.ts'
import {AspectRatio, Image, Menu, Stack, Text} from '@mantine/core'
import albumPlaceholder from '../../assets/image-placeholder-1.jpg'
import { useState } from 'react'
import { useAppDispatch } from '../../state/store.ts'
import { openArtistDrawer } from '../../state/globalSlice.ts'
import { useNavigate } from 'react-router-dom'
import useContextMenu from "../../hooks/useContextMenu.ts";
import {IconTrash} from "@tabler/icons-react";
import {toast} from "react-toastify";
import {useDeleteAlbumMutation} from "../../state/albumsApi.ts";

interface AlbumCardProps {
  album: Album
}

function AlbumCard({ album }: AlbumCardProps) {
  const dispatch = useAppDispatch()
  const navigate = useNavigate()

  const [deleteAlbumMutation] = useDeleteAlbumMutation()

  const [isImageHovered, setIsImageHovered] = useState(false)
  const [openedMenu, menuDropdownProps, { openMenu, onMenuChange }] = useContextMenu()

  function handleClick() {
    navigate(`/album/${album.id}`)
  }

  function handleArtistClick(artistId: string) {
    dispatch(openArtistDrawer(artistId))
  }

  function handleDelete() {
    deleteAlbumMutation(album.id)
    toast.success(`${album.title} deleted!`)
  }

  return (
    <Stack
      align={'center'}
      gap={0}
      style={{ transition: '0.3s', ...(isImageHovered && { transform: 'scale(1.1)' }) }}
      w={150}
    >
      <Menu shadow={'lg'} opened={openedMenu} onChange={onMenuChange}>
        <Menu.Target>
          <AspectRatio>
            <Image
              onMouseEnter={() => setIsImageHovered(true)}
              onMouseLeave={() => setIsImageHovered(false)}
              radius={'lg'}
              src={album.imageUrl}
              fallbackSrc={albumPlaceholder}
              onClick={handleClick}
              onContextMenu={openMenu}
              sx={(theme) => ({
                cursor: 'pointer',
                transition: '0.3s',
                boxShadow: theme.shadows.xxl,
                '&:hover': {
                  boxShadow: theme.shadows.xxl_hover
                }
              })}
            />
          </AspectRatio>
        </Menu.Target>

        <Menu.Dropdown {...menuDropdownProps}>
          <Menu.Item c={'red'} leftSection={<IconTrash size={14} />} onClick={handleDelete}>
            Delete
          </Menu.Item>
        </Menu.Dropdown>
      </Menu>

      <Stack w={'100%'} align={'center'} pt={'xs'} gap={0} style={{ overflow: 'hidden' }}>
        <Text fw={600} lineClamp={2} ta={'center'}>
          {album.title}
        </Text>
        {album.artist ? (
          <Text
            fw={500}
            ta={'center'}
            c={'dimmed'}
            truncate={'end'}
            onClick={() => handleArtistClick(album.artist.id)}
            sx={{
              cursor: 'pointer',
              '&:hover': {
                textDecoration: 'underline'
              }
            }}
          >
            {album.artist.name}
          </Text>
        ) : (
          <Text c={'dimmed'} ta={'center'} fs={'oblique'}>
            Unknown
          </Text>
        )}
      </Stack>
    </Stack>
  )
}

export default AlbumCard
