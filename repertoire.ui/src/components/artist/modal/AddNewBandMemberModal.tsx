import { Button, Group, Modal, Stack, TextInput } from '@mantine/core'
import { useState } from 'react'
import { FileWithPath } from '@mantine/dropzone'
import { useForm } from '@mantine/form'
import { zod4Resolver } from 'mantine-form-zod-resolver'
import { AddNewBandMemberForm, addNewBandMemberSchema } from '../../../validation/artistsForm.ts'
import { toast } from 'react-toastify'
import ImageDropzoneWithPreview from '../../@ui/image/ImageDropzoneWithPreview.tsx'
import {
  useCreateBandMemberMutation,
  useSaveImageToBandMemberMutation
} from '../../../state/api/artistsApi.ts'
import { IconUserFilled } from '@tabler/icons-react'
import BandMemberRoleMultiSelect from '../../@ui/form/select/multi/BandMemberRoleMultiSelect.tsx'
import ColorInputButton from '../../@ui/form/input/button/ColorInputButton.tsx'
import bandMemberColorSwatches from '../../../data/artist/bandMemberColorSwatches.ts'
import { useDidUpdate } from '@mantine/hooks'

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
  useDidUpdate(() => {
    form.setFieldValue('roleIds', roleIds)
    form.validateField('roleIds')
  }, [roleIds])

  const onCloseWithImage = () => {
    onClose()
    setImage(null)
  }

  const form = useForm<AddNewBandMemberForm>({
    mode: 'uncontrolled',
    initialValues: {
      name: '',
      roleIds: []
    },
    validateInputOnBlur: true,
    validateInputOnChange: false,
    clearInputErrorOnChange: true,
    validate: zod4Resolver(addNewBandMemberSchema)
  })

  async function addBandMember({ name, roleIds }: AddNewBandMemberForm) {
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
                  setIds={setRoleIds}
                  label={'Roles'}
                  placeholder={'Select roles'}
                  withAsterisk
                  pr={'lg'}
                  error={form.getInputProps('roleIds').error}
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
