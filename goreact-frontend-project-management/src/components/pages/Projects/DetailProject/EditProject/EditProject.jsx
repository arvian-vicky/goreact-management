import { GroupAdd, Remove, RemoveCircle, Settings } from '@mui/icons-material';
import { IconButton, Menu, MenuItem } from '@mui/material';
import { useState } from 'react';
import ModalUpdateDeadline from './Modals/ModalUpdateDeadline';
import ModalCloseProject from './Modals/ModalCloseProject';
import Dropdown from '@/components/ui/Dropdown';
import boards from '@/services/api/boards';
import { useNavigate, useParams } from 'react-router';
import { CalendarIcon } from '@mui/x-date-pickers';
import ModalAddNewMember from './Modals/ModalAddNewMember';

const EditProjectButton = ({ data }) => {
  const { id } = useParams();
  const navigate = useNavigate();

  const [anchorEl, setAnchorEl] = useState(null);
  const openMenu = Boolean(anchorEl);

  const [openUpdateDeadline, setOpenUpdateDeadline] = useState(false);
  const [openCloseProject, setOpenCloseProject] = useState(false);
  const [openModalAddNewMember, setOpenModalAddNewMember] = useState(false);

  const handleClick = (event) => {
    setAnchorEl(event.currentTarget);
  };
  const handleCloseMenu = () => {
    setAnchorEl(null);
  };

  const handleOpenUpdateDeadline = () => {
    setOpenUpdateDeadline(true);
    handleCloseMenu();
  };

  const handleOpenCloseProject = () => {
    setOpenCloseProject(true);
    handleCloseMenu();
  };

  const handleCloseUpdateDeadline = () => {
    setOpenUpdateDeadline(false);
  };

  const handleCloseProject = async () => {
    await boards.remove(id);
    setOpenCloseProject(false);
    navigate('/projects');
  };

  const handleOpenModalAddNewMember = () => {
    setOpenModalAddNewMember(true);
    handleCloseMenu();
  };

  const handleCloseModalAddNewMember = () => {
    setOpenModalAddNewMember(false);
  };

  return (
    <>
      <Dropdown
        icon={<Settings />}
        options={[
          {
            icon: <RemoveCircle />,
            label: 'Tutup proyek ini',
            onClick: handleOpenCloseProject,
          },
          {
            icon: <CalendarIcon />,
            label: 'Ubah deadline',
            onClick: handleOpenUpdateDeadline,
          },
          {
            icon: <GroupAdd />,
            label: 'Tambah anggota tim',
            onClick: handleOpenModalAddNewMember,
          },
        ]}
      />

      <ModalUpdateDeadline
        data={data}
        open={openUpdateDeadline}
        handleClose={handleCloseUpdateDeadline}
      />
      <ModalCloseProject
        open={openCloseProject}
        handleClose={handleCloseProject}
      />
      <ModalAddNewMember
        open={openModalAddNewMember}
        handleClose={handleCloseModalAddNewMember}
      />
    </>
  );
};

export default EditProjectButton;
