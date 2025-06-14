import { Button, Group, LoadingOverlay, Modal, Stack, TextInput, Tooltip } from '@mantine/core'
import { useEffect, useState } from 'react'
import { FileWithPath } from '@mantine/dropzone'
import { useForm } from '@mantine/form'
import { zod4Resolver } from 'mantine-form-zod-resolver'
import { EditBandMemberForm, editBandMemberSchema } from '../../../validation/artistsForm.ts'
import { toast } from 'react-toastify'
import ImageDropzoneWithPreview from '../../@ui/image/ImageDropzoneWithPreview.tsx'
import {
  useDeleteImageFromBandMemberMutation,
  useSaveImageToBandMemberMutation,
  useUpdateBandMemberMutation
} from '../../../state/api/artistsApi.ts'
import { IconUserFilled } from '@tabler/icons-react'
import BandMemberRoleMultiSelect from '../../@ui/form/select/multi/BandMemberRoleMultiSelect.tsx'
import { useDidUpdate } from '@mantine/hooks'
import { BandMember } from '../../../types/models/Artist.ts'
import ColorInputButton from '../../@ui/form/input/button/ColorInputButton.tsx'
import bandMemberColorSwatches from '../../../data/artist/bandMemberColorSwatches.ts'

interface EditBandMemberModalProps {
  opened: boolean
  onClose: () => void
  bandMember: BandMember
}

function EditBandMemberModal({ opened, onClose, bandMember }: EditBandMemberModalProps) {
  const [updateBandMemberMutation, { isLoading: isUpdateLoading }] = useUpdateBandMemberMutation()
  const [saveImageMutation, { isLoading: isSaveImageLoading }] = useSaveImageToBandMemberMutation()
  const [deleteImageMutation, { isLoading: isDeleteImageLoading }] =
    useDeleteImageFromBandMemberMutation()
  const isLoading = isUpdateLoading || isSaveImageLoading || isDeleteImageLoading

  const [memberHasChanged, setMemberHasChanged] = useState(false)
  const [imageHasChanged, setImageHasChanged] = useState(false)
  const hasChanged = memberHasChanged || imageHasChanged

  function areRolesEqual() {
    if (roleIds.length !== bandMember.roles.length) return false

    for (const role of bandMember.roles) {
      if (!roleIds.includes(role.id)) return false
    }

    return true
  }

  const form = useForm<EditBandMemberForm>({
    mode: 'uncontrolled',
    initialValues: {
      name: bandMember.name,
      color: bandMember.color,
      image: bandMember.imageUrl,
      roleIds: bandMember.roles.map((role) => role.id)
    },
    validateInputOnBlur: true,
    validateInputOnChange: false,
    clearInputErrorOnChange: true,
    validate: zod4Resolver(editBandMemberSchema),
    onValuesChange: (values) => {
      setMemberHasChanged(
        values.name !== bandMember.name || values.color !== bandMember.color || !areRolesEqual()
      )
      setImageHasChanged(values.image !== bandMember.imageUrl)
    }
  })

  const [color, setColor] = useState<string>(bandMember.color)
  useEffect(() => form.setFieldValue('color', color), [color])

  const [roleIds, setRoleIds] = useState<string[]>(bandMember.roles.map((r) => r.id))
  useDidUpdate(() => {
    form.setFieldValue('roleIds', roleIds)
    if (opened) form.validateField('roleIds')
  }, [roleIds])

  const [image, setImage] = useState<FileWithPath | string>(bandMember.imageUrl)
  useEffect(() => form.setFieldValue('image', image), [image])
  useDidUpdate(() => setImage(bandMember.imageUrl), [bandMember])

  async function addBandMember({ name, color, image, roleIds }: EditBandMemberForm) {
    if (memberHasChanged)
      await updateBandMemberMutation({
        id: bandMember.id,
        name: name.trim(),
        color: color,
        roleIds: roleIds
      }).unwrap()

    if (image && typeof image !== 'string')
      await saveImageMutation({
        id: bandMember.id,
        image: image as FileWithPath
      })
    else if (!image && bandMember.imageUrl) await deleteImageMutation(bandMember.id)

    toast.info(`${name} updated!`)
    onClose()
    setMemberHasChanged(false)
    setImageHasChanged(false)
  }

  return (
    <Modal opened={opened} onClose={onClose} title={'Edit Band Member'}>
      <form onSubmit={form.onSubmit(addBandMember)}>
        <Stack py={'xs'} pl={'xs'} pr={0}>
          <LoadingOverlay visible={isLoading} loaderProps={{ type: 'bars' }} />

          <Group>
            <ImageDropzoneWithPreview
              image={image}
              setImage={setImage}
              defaultValue={bandMember.imageUrl}
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

          <Tooltip
            disabled={hasChanged}
            label={'You need to make a change before saving'}
            position="bottom"
          >
            <Button
              type={'submit'}
              data-disabled={!hasChanged}
              onClick={(e) => (!hasChanged ? e.preventDefault() : {})}
            >
              Save
            </Button>
          </Tooltip>
        </Stack>
      </form>
    </Modal>
  )
}

export default EditBandMemberModal
