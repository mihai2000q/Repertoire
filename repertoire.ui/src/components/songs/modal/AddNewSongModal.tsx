import { Button, ComboboxItem, Group, Modal, Space, Stepper, Text } from '@mantine/core'
import { useForm } from '@mantine/form'
import { zod4Resolver } from 'mantine-form-zod-resolver'
import { AddNewSongForm, addNewSongSchema } from '../../../validation/songsForm.ts'
import { useCreateSongMutation, useSaveImageToSongMutation } from '../../../state/api/songsApi.ts'
import { useState } from 'react'
import { toast } from 'react-toastify'
import { FileWithPath } from '@mantine/dropzone'
import AddNewSongModalFirstStep from './AddNewSongModalFirstStep.tsx'
import AddNewSongModalSecondStep from './AddNewSongModalSecondStep.tsx'
import AddNewSongModalFinalStep from './AddNewSongModalFinalStep.tsx'
import { useDidUpdate, useListState } from '@mantine/hooks'
import { AlbumSearch, ArtistSearch } from '../../../types/models/Search.ts'

export interface AddNewSongModalSongSection {
  id: string
  name: string
  type: ComboboxItem | null
  errors: { property: string }[]
}

interface AddNewSongModalProps {
  opened: boolean
  onClose: () => void
}

function AddNewSongModal({ opened, onClose }: AddNewSongModalProps) {
  const [createSongMutation, { isLoading: isCreateSongLoading }] = useCreateSongMutation()
  const [saveImageMutation, { isLoading: isSaveImageLoading }] = useSaveImageToSongMutation()
  const isLoading = isCreateSongLoading || isSaveImageLoading

  const [artist, setArtist] = useState<ArtistSearch>(null)
  const [album, setAlbum] = useState<AlbumSearch>(null)

  const [guitarTuning, setGuitarTuning] = useState<ComboboxItem>(null)
  const [difficulty, setDifficulty] = useState<ComboboxItem>(null)
  const [sections, sectionsHandlers] = useListState<AddNewSongModalSongSection>([])

  const [image, setImage] = useState<FileWithPath>(null)

  const onCloseWithImage = () => {
    onClose()
    setImage(null)
  }

  const [isArtistDisabled, setIsArtistDisabled] = useState(false)
  useDidUpdate(() => setIsArtistDisabled(album !== null), [album])

  const form = useForm<AddNewSongForm>({
    mode: 'uncontrolled',
    initialValues: {
      title: '',
      description: ''
    },
    validateInputOnBlur: true,
    validateInputOnChange: false,
    clearInputErrorOnChange: true,
    validate: zod4Resolver(addNewSongSchema),
    enhanceGetInputProps: (payload) => ({
      disabled: isArtistDisabled && payload.field === 'artistName'
    })
  })

  const [activeStep, setActiveStep] = useState(0)
  const handleActiveStepChange = (activeStep: ((prevState: number) => number) | number) => {
    const localSections = sections.map((s) => ({
      ...s,
      errors: [
        ...(s.name.trim() === '' ? [{ property: 'name' }] : []),
        ...(!s.type ? [{ property: 'type' }] : [])
      ]
    }))
    sectionsHandlers.setState(localSections)
    if (form.validate().hasErrors || localSections.some((s) => s.errors.length > 0)) return
    setActiveStep(activeStep)
  }
  const prevStep = () => handleActiveStepChange((current) => current - 1)
  const nextStep = () => handleActiveStepChange((current) => current + 1)

  async function addSong({
    title,
    description,
    artistName,
    albumTitle,
    bpm,
    releaseDate,
    songsterrLink,
    youtubeLink
  }: AddNewSongForm) {
    title = title.trim()
    artistName = artistName?.trim() === '' ? undefined : artistName?.trim()
    albumTitle = albumTitle?.trim() === '' ? undefined : albumTitle?.trim()
    songsterrLink = songsterrLink?.trim() === '' ? undefined : songsterrLink?.trim()
    youtubeLink = youtubeLink?.trim() === '' ? undefined : youtubeLink?.trim()
    bpm = typeof bpm === 'string' ? undefined : bpm

    const res = await createSongMutation({
      title: title,
      description: description,
      bpm: bpm,
      difficulty: difficulty?.value,
      releaseDate: releaseDate,
      songsterrLink: songsterrLink,
      youtubeLink: youtubeLink,
      sections: sections.map((s) => ({ name: s.name.trim(), typeId: s.type.value })),
      guitarTuningId: guitarTuning?.value,
      albumId: album?.id,
      artistId: album ? undefined : artist?.id,
      artistName: artist ? undefined : artistName,
      albumTitle: album ? undefined : albumTitle
    }).unwrap()

    if (image) await saveImageMutation({ image: image, id: res.id }).unwrap()

    toast.success(`${title} added!`)

    onCloseWithImage()
    setActiveStep(0)
    form.reset()
    setArtist(null)
    setAlbum(null)
    setGuitarTuning(null)
    setDifficulty(null)
    setIsArtistDisabled(false)
    sectionsHandlers.setState([])
  }

  return (
    <Modal opened={opened} onClose={onCloseWithImage} title={'Add New Song'} size={500}>
      <Modal.Body p={'xs'}>
        <form onSubmit={form.onSubmit(addSong)}>
          <Stepper iconSize={35} active={activeStep} onStepClick={handleActiveStepChange}>
            <Stepper.Step
              label={
                <Group gap={0}>
                  First Step
                  <Text inline pl={1} c={'red'}>
                    *
                  </Text>
                </Group>
              }
              description={'Basic Info'}
            >
              <AddNewSongModalFirstStep
                form={form}
                artist={artist}
                setArtist={setArtist}
                album={album}
                setAlbum={setAlbum}
              />
            </Stepper.Step>

            <Stepper.Step label={'Second Step'} description={'More Info'}>
              <AddNewSongModalSecondStep
                form={form}
                sections={sections}
                sectionsHandlers={sectionsHandlers}
                guitarTuning={guitarTuning}
                setGuitarTuning={setGuitarTuning}
                difficulty={difficulty}
                setDifficulty={setDifficulty}
                album={album}
              />
            </Stepper.Step>

            <Stepper.Step label={'Final Step'} description={'Web & Media'}>
              <AddNewSongModalFinalStep
                form={form}
                image={image}
                setImage={setImage}
                album={album}
              />
            </Stepper.Step>
          </Stepper>

          <Group pt={'xs'} gap={'xs'}>
            <Space flex={1} />
            {activeStep !== 0 && (
              <Button variant={'subtle'} onClick={prevStep}>
                Previous
              </Button>
            )}
            {activeStep !== 2 && <Button onClick={nextStep}>Next</Button>}
            {activeStep === 2 && (
              <Button type={'submit'} loading={isLoading}>
                Submit
              </Button>
            )}
          </Group>
        </form>
      </Modal.Body>
    </Modal>
  )
}

export default AddNewSongModal
