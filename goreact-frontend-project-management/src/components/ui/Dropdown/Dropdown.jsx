import { Box, IconButton, Menu, MenuItem, Typography } from '@mui/material';
import { useState } from 'react';

const Dropdown = ({ options, icon }) => {
  const [anchorEl, setAnchorEl] = useState(null);
  const open = Boolean(anchorEl);
  const handleClick = (e) => {
    setAnchorEl(e.currentTarget);
  };
  const handleClose = () => {
    setAnchorEl(null);
  };
  return (
    <>
      <IconButton onClick={handleClick}>{icon}</IconButton>
      <Menu anchorEl={anchorEl} open={open} onClose={handleClose}>
        {options?.map((option, index) => (
          <MenuItem
            key={index}
            onClick={() => {
              option.onClick();
              handleClose();
            }}
            sx={{
              display: 'flex',
              justifyContent: 'space-between',
              alignItems: 'center',
              gap: 3,
            }}
          >
            <Typography variant={'body1'}>{option.label}</Typography>
            <Box>{option.icon ?? <></>}</Box>
          </MenuItem>
        ))}
      </Menu>
    </>
  );
};

export default Dropdown;
