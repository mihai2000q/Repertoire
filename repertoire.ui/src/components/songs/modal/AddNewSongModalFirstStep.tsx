import { Box, Group, Stack, Textarea, TextInput, Tooltip } from '@mantine/core'
import { UseFormReturnType } from '@mantine/form'
import ArtistAutocomplete from '../../@ui/form/input/ArtistAutocomplete.tsx'
import Artist from '../../../types/models/Artist.ts'
import AlbumAutocomplete from '../../@ui/form/input/AlbumAutocomplete.tsx'
import Album from '../../../types/models/Album.ts'
import { AddNewSongForm } from '../../../validation/songsForm.ts'
import { IconInfoCircleFilled } from '@tabler/icons-react'
import { useDidUpdate } from '@mantine/hooks'

interface AddNewSongModalFirstStepProps {
  form: UseFormReturnType<AddNewSongForm>
  artist: Artist
  setArtist: (artist: Artist) => void
  album: Album
  setAlbum: (album: Album) => void
}

function AddNewSongModalFirstStep({
  form,
  artist,
  setArtist,
  album,
  setAlbum
}: AddNewSongModalFirstStepProps) {
  useDidUpdate(() => {
    setArtist(album?.artist)
    form.setFieldValue('artistName', album?.artist?.name)
  }, [album])

  return (
    <Stack>
      <TextInput
        withAsterisk={true}
        maxLength={100}
        label="Title"
        placeholder="The title of the song"
        key={form.key('title')}
        {...form.getInputProps('title')}
      />

      <Group>
        <AlbumAutocomplete
          album={album}
          setAlbum={setAlbum}
          key={form.key('albumTitle')}
          setValue={(v) => form.setFieldValue('albumTitle', v)}
          {...form.getInputProps('albumTitle')}
        />

        <Group flex={1} gap={0}>
          <ArtistAutocomplete
            artist={artist}
            setArtist={setArtist}
            key={form.key('artistName')}
            setValue={(v) => form.setFieldValue('artistName', v)}
            {...form.getInputProps('artistName')}
          />
          {album && (
            <Box c={'primary.8'} mt={'lg'} ml={4}>
              <Tooltip
                multiline
                w={210}
                ta={'center'}
                label={'Song will inherit artist from album (even if it has one or not)'}
              >
                <IconInfoCircleFilled aria-label={'artist-info'} size={18} />
              </Tooltip>
            </Box>
          )}
        </Group>
      </Group>

      <Textarea
        label="Description"
        placeholder="Enter Description"
        autosize={true}
        minRows={3}
        maxRows={6}
        key={form.key('description')}
        {...form.getInputProps('description')}
      />
    </Stack>
  )
}

export default AddNewSongModalFirstStep
