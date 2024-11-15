import { Button, ComboboxItem, Group, Modal, Space, Stepper, Text } from '@mantine/core'
import { useForm, zodResolver } from '@mantine/form'
import { AddNewSongForm, addNewSongValidation } from '../../../validation/songsForm.ts'
import { useCreateSongMutation, useSaveImageToSongMutation } from '../../../state/songsApi.ts'
import { useState } from 'react'
import { toast } from 'react-toastify'
import { FileWithPath } from '@mantine/dropzone'
import AddNewSongModalFirstStep from './AddNewSongModalFirstStep.tsx'
import AddNewSongModalSecondStep from './AddNewSongModalSecondStep.tsx'
import AddNewSongModelFinalStep from './AddNewSongModelFinalStep.tsx'
import { useListState } from '@mantine/hooks'

export interface AddNewSongModalSongSection {
  id: string
  name: string
  type: ComboboxItem | null
}

interface AddNewSongModalProps {
  opened: boolean
  onClose: () => void
}

function AddNewSongModal({ opened, onClose }: AddNewSongModalProps) {
  const [createSongMutation, { isLoading: isCreateSongLoading }] = useCreateSongMutation()
  const [saveImageMutation, { isLoading: isSaveImageLoading }] = useSaveImageToSongMutation()
  const isLoading = isCreateSongLoading || isSaveImageLoading

  const [guitarTuning, setGuitarTuning] = useState<ComboboxItem>(null)
  const [difficulty, setDifficulty] = useState<ComboboxItem>(null)
  const [sections, sectionsHandlers] = useListState<AddNewSongModalSongSection>([])

  const [image, setImage] = useState<FileWithPath>(null)

  const onCloseWithImage = () => {
    onClose()
    setImage(null)
  }

  const form = useForm({
    mode: 'uncontrolled',
    initialValues: {
      title: '',
      description: ''
    } as AddNewSongForm,
    validateInputOnBlur: true,
    validateInputOnChange: false,
    clearInputErrorOnChange: true,
    validate: zodResolver(addNewSongValidation)
  })

  const [activeStep, setActiveStep] = useState(0)
  const handleActiveStepChange = (activeStep: ((prevState: number) => number) | number) => {
    if (form.validate().hasErrors) return
    setActiveStep(activeStep)
  }
  const prevStep = () => handleActiveStepChange((current) => current - 1)
  const nextStep = () => handleActiveStepChange((current) => current + 1)

  async function addSong({
    title,
    description,
    bpm,
    releaseDate,
    songsterrLink,
    youtubeLink
  }: AddNewSongForm) {
    title = title.trim()
    songsterrLink = songsterrLink?.trim()
    youtubeLink = youtubeLink?.trim()

    const res = await createSongMutation({
      title: title,
      description: description,
      bpm: bpm,
      difficulty: difficulty?.value,
      releaseDate: releaseDate,
      songsterrLink: songsterrLink,
      youtubeLink: youtubeLink,
      sections: sections.map((s) => ({ name: s.name, typeId: s.type.value })),
      guitarTuningId: guitarTuning?.value
    }).unwrap()

    if (image) await saveImageMutation({ image: image, id: res.id }).unwrap()

    toast.success(`${title} added!`)

    onCloseWithImage()
    setActiveStep(0)
    form.reset()
    setGuitarTuning(null)
    setDifficulty(null)
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
              <AddNewSongModalFirstStep form={form} />
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
              />
            </Stepper.Step>

            <Stepper.Step label={'Final Step'} description={'Web & Media'}>
              <AddNewSongModelFinalStep form={form} image={image} setImage={setImage} />
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
              <Button type={'submit'} disabled={isLoading}>
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
