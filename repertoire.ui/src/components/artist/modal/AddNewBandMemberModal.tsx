import { Button, Group, Modal, Stack, TextInput } from '@mantine/core'
import { useState } from 'react'
import { FileWithPath } from '@mantine/dropzone'
import { useForm, zodResolver } from '@mantine/form'
import {
  AddNewBandMemberForm,
  addNewBandMemberValidation
} from '../../../validation/artistsForm.ts'
import { toast } from 'react-toastify'
import ImageDropzoneWithPreview from '../../@ui/image/ImageDropzoneWithPreview.tsx'
import {
  useCreateBandMemberMutation,
  useSaveImageToBandMemberMutation
} from '../../../state/api/artistsApi.ts'
import { IconUserFilled } from '@tabler/icons-react'
import BandMemberRoleMultiSelect from '../../@ui/form/select/multi/BandMemberRoleMultiSelect.tsx'
import ColorInputButton from '../../@ui/form/input/ColorInputButton.tsx'
import bandMemberColorSwatches from '../../../data/artist/bandMemberColorSwatches.ts'

interface AddNewBandMemberModalProps {
  opened: boolean
  onClose: () => void
  artistId: string
}

function AddNewBandMemberModal({ opened, onClose, artistId }: AddNewBandMemberModalProps) {
  const [createBandMemberMutation, { isLoading: isCreateLoading }] = useCreateBandMemberMutation()
  const [saveImageMutation, { isLoading: isSaveImageLoading }] = useSaveImageToBandMemberMutation()
  const isLoading = isCreateLoading || isSaveImageLoading

  const [image, setImage] = useState<FileWithPath>(null)
  const [color, setColor] = useState<string>()
  const [roleIds, setRoleIds] = useState<string[]>([])
  const [rolesError, setRolesError] = useState<boolean>(false)

  const onCloseWithImage = () => {
    onClose()
    setImage(null)
    setRolesError(false)
  }

  const form = useForm({
    mode: 'uncontrolled',
    initialValues: {
      name: ''
    } as AddNewBandMemberForm,
    validateInputOnBlur: true,
    validateInputOnChange: false,
    clearInputErrorOnChange: true,
    validate: zodResolver(addNewBandMemberValidation)
  })

  async function addBandMember({ name }: AddNewBandMemberForm) {
    if (roleIds.length === 0) {
      setRolesError(true)
      return
    }

    name = name.trim()

    const res = await createBandMemberMutation({
      name: name,
      color: color,
      roleIds: roleIds,
      artistId: artistId
    }).unwrap()

    if (image) await saveImageMutation({ image: image, id: res.id }).unwrap()

    toast.success(`${name} added!`)
    onCloseWithImage()
    setColor(undefined)
    setRoleIds([])
    form.reset()
  }

  return (
    <Modal opened={opened} onClose={onCloseWithImage} title={'Add New Band Member'}>
      <Modal.Body py={'xs'} pl={'xs'} pr={0}>
        <form onSubmit={form.onSubmit(addBandMember)}>
          <Stack>
            <Group>
              <ImageDropzoneWithPreview
                image={image}
                setImage={setImage}
                radius={'50%'}
                w={100}
                h={100}
                icon={<IconUserFilled size={45} />}
              />

              <Stack flex={1}>
                <Group align={'start'}>
                  <TextInput
                    flex={1}
                    withAsterisk={true}
                    maxLength={100}
                    label="Name"
                    placeholder="Name of band member"
                    key={form.key('name')}
                    {...form.getInputProps('name')}
                  />

                  <ColorInputButton
                    color={color}
                    setColor={setColor}
                    swatches={bandMemberColorSwatches}
                  />
                </Group>

                <BandMemberRoleMultiSelect
                  ids={roleIds}
                  setIds={(ids) => {
                    setRoleIds(ids)
                    setRolesError(ids.length === 0)
                  }}
                  label={'Roles'}
                  placeholder={'Select roles'}
                  withAsterisk
                  pr={'lg'}
                  error={rolesError && 'Please select at least one role'}
                />
              </Stack>
            </Group>

            <Button style={{ alignSelf: 'center' }} type={'submit'} loading={isLoading}>
              Submit
            </Button>
          </Stack>
        </form>
      </Modal.Body>
    </Modal>
  )
}

export default AddNewBandMemberModal
