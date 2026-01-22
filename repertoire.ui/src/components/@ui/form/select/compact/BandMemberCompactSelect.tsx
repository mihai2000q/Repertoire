import {
  ActionIcon,
  alpha,
  Avatar,
  Combobox,
  ComboboxProps,
  Group,
  ScrollArea,
  Stack,
  Text,
  Tooltip,
  UnstyledButton,
  useCombobox
} from '@mantine/core'
import { BandMember } from '../../../../../types/models/Artist.ts'
import { IconSearch, IconUser } from '@tabler/icons-react'
import { RefObject, useEffect, useState } from 'react'
import { useInputState } from '@mantine/hooks'

interface BandMemberCompactSelectProps extends ComboboxProps {
  bandMember: BandMember | null
  setBandMember: (bandMember: BandMember | null) => void
  bandMembers: BandMember[] | undefined
  ref: RefObject<HTMLButtonElement>
  tooltipLabel?: string
}

function BandMemberCompactSelect({
  bandMember,
  setBandMember,
  bandMembers,
  ref,
  tooltipLabel,
  ...others
}: BandMemberCompactSelectProps) {
  const [value, setValue] = useState<string>(bandMember?.name ?? '')
  const [search, setSearch] = useInputState(bandMember?.name ?? '')
  useEffect(() => {
    setValue(bandMember?.name ?? '')
    setSearch(bandMember?.name ?? '')
  }, [bandMember])

  const combobox = useCombobox({
    onDropdownClose: () => {
      combobox.resetSelectedOption()
      combobox.focusTarget()
      setSearch(value)
    }
  })

  const filteredMembers =
    search.trim() !== ''
      ? bandMembers.filter((member) =>
          member.name.toLowerCase().includes(search.toLowerCase().trim())
        )
      : bandMembers

  const BandMemberOption = ({ member }: { member: BandMember }) => (
    <Combobox.Option
      key={member.id}
      value={member.name}
      aria-label={member.name}
      onClick={() => setBandMember(bandMember === member ? null : member)}
    >
      <Group gap={'xs'} wrap={'nowrap'}>
        <Avatar
          size={'sm'}
          variant={'light'}
          color={member.color}
          src={member.imageUrl}
          alt={member.imageUrl && member.name}
        >
          <IconUser size={14} />
        </Avatar>
        <Text lh={'xxs'} fw={500} lineClamp={2}>
          {member.name}
        </Text>
      </Group>
    </Combobox.Option>
  )

  function handleSubmit(valueString: string) {
    setValue(valueString)
    setSearch(valueString)
    combobox.closeDropdown()
  }

  function handleClear() {
    setBandMember(null)
  }

  return (
    <Combobox
      onOptionSubmit={handleSubmit}
      store={combobox}
      withArrow
      onEnterTransitionEnd={combobox.focusSearchInput}
      {...others}
    >
      <Combobox.Target withAriaAttributes={false}>
        {bandMember ? (
          <Tooltip
            label={
              <Text fz={'sm'} c={'white'} lineClamp={2}>
                {bandMember.name} is selected
              </Text>
            }
            openDelay={200}
            multiline={true}
            maw={250}
            disabled={combobox.dropdownOpened}
          >
            <UnstyledButton
              ref={ref}
              aria-label={bandMember.name}
              onClick={() => combobox.toggleDropdown()}
              style={{ borderRadius: '50%' }}
            >
              <Avatar
                size={28}
                color={bandMember.color}
                src={bandMember.imageUrl}
                sx={{ transition: '0.25s', '&:hover': { filter: 'brightness(0.7)' } }}
              >
                <IconUser size={15} />
              </Avatar>
            </UnstyledButton>
          </Tooltip>
        ) : (
          <Tooltip
            label={
              bandMembers
                ? (tooltipLabel ?? 'Choose a band member')
                : 'The artist of the song is either not set or not a band'
            }
            ta={'center'}
            openDelay={bandMembers ? 200 : 0}
            maw={200}
            multiline
            position={'top'}
            disabled={combobox.dropdownOpened}
          >
            <ActionIcon
              ref={ref}
              aria-label={'select-band-member'}
              radius={'50%'}
              sx={(theme) => ({
                transition: '0.25s',
                backgroundColor: 'transparent',
                color: theme.colors.gray[5],
                border: `1px dashed ${theme.colors.gray[5]}`,

                '&:hover': {
                  backgroundColor: theme.colors.gray[1],
                  color: alpha(theme.colors.gray[6], 0.85),
                  borderColor: alpha(theme.colors.gray[6], 0.85)
                },
                '&[data-disabled="true"]': {
                  backgroundColor: 'transparent',
                  color: alpha(theme.colors.gray[4], 0.8),
                  borderColor: alpha(theme.colors.gray[4], 0.8)
                }
              })}
              onClick={() => combobox.toggleDropdown()}
              disabled={bandMembers === undefined}
            >
              <IconUser size={15} />
            </ActionIcon>
          </Tooltip>
        )}
      </Combobox.Target>

      <Combobox.Dropdown miw={180} pt={2} px={'xxs'} pb={0}>
        <Combobox.Search
          size={'xs'}
          px={'xxs'}
          pos={'relative'}
          maxLength={100}
          aria-label={'search'}
          placeholder={'Search'}
          leftSection={<IconSearch size={12} />}
          rightSection={bandMember && <Combobox.ClearButton onClear={handleClear} />}
          value={search}
          onChange={setSearch}
          sx={{
            '.mantine-Input-section': {
              position: 'absolute',
              pointerEvents: 'all'
            },
            '.mantine-Combobox-input': { paddingTop: 1 }
          }}
        />
        <Combobox.Options>
          <ScrollArea.Autosize mah={200} scrollbarSize={5}>
            <Stack gap={0} pb={'xxs'}>
              {bandMembers?.length === 0 ? (
                <Combobox.Empty>Artist has no members</Combobox.Empty>
              ) : filteredMembers?.length === 0 ? (
                <Combobox.Empty>No members found</Combobox.Empty>
              ) : (
                filteredMembers?.map((bandMember) => (
                  <BandMemberOption key={bandMember.id} member={bandMember} />
                ))
              )}
            </Stack>
          </ScrollArea.Autosize>
        </Combobox.Options>
      </Combobox.Dropdown>
    </Combobox>
  )
}

export default BandMemberCompactSelect
