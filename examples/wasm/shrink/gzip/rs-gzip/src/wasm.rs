use std::sync::RwLock;

use crate::encode2gzip::enc2gz::CompressionLevel;

static INPUT: RwLock<Vec<u8>> = RwLock::new(vec![]);
static OUTPUT: RwLock<Vec<u8>> = RwLock::new(vec![]);

pub fn set_size(i: &mut Vec<u8>, sz: usize, val: u8) -> usize {
    i.resize(sz, val);
    i.capacity()
}

pub fn _set_input_size(sz: usize, val: u8) -> Result<usize, &'static str> {
    let mut guard = INPUT
        .try_write()
        .map_err(|_| "unable to write lock the input")?;
    let m: &mut Vec<u8> = &mut guard;
    let cap: usize = set_size(m, sz, val);
    Ok(cap)
}

pub fn _set_output_size(sz: usize, val: u8) -> Result<usize, &'static str> {
    let mut guard = OUTPUT
        .try_write()
        .map_err(|_| "unable to write lock the input")?;
    let m: &mut Vec<u8> = &mut guard;
    let cap: usize = set_size(m, sz, val);
    Ok(cap)
}

#[allow(unsafe_code)]
#[no_mangle]
pub extern "C" fn set_input_size(size: u32) -> i32 {
    size.try_into()
        .ok()
        .and_then(|u: usize| _set_input_size(u, 0).ok())
        .and_then(|cap: usize| cap.try_into().ok())
        .unwrap_or(-1)
}

#[allow(unsafe_code)]
#[no_mangle]
pub extern "C" fn set_output_size(size: u32) -> i32 {
    size.try_into()
        .ok()
        .and_then(|u: usize| _set_output_size(u, 0).ok())
        .and_then(|cap: usize| cap.try_into().ok())
        .unwrap_or(-1)
}

pub fn _get_input_size() -> Result<usize, &'static str> {
    let guard = INPUT
        .try_read()
        .map_err(|_| "unable to read lock the input")?;
    let s: &[u8] = &guard;
    Ok(s.len())
}

#[allow(unsafe_code)]
#[no_mangle]
pub extern "C" fn get_input_size() -> i32 {
    _get_input_size()
        .ok()
        .and_then(|u| u.try_into().ok())
        .unwrap_or(-1)
}

#[allow(unsafe_code)]
#[no_mangle]
pub extern "C" fn get_out_estimate() -> i32 {
    let isz: i32 = get_input_size();
    match isz {
        0.. => isz << 1,
        _ => -1,
    }
}

pub fn offset(s: &[u8]) -> *const u8 {
    s.as_ptr()
}

pub fn _offset_i() -> Result<*const u8, &'static str> {
    let guard = INPUT
        .try_read()
        .map_err(|_| "unable to read lock the input")?;
    let s: &[u8] = &guard;
    Ok(offset(s))
}

pub fn _offset_o() -> Result<*const u8, &'static str> {
    let guard = OUTPUT
        .try_read()
        .map_err(|_| "unable to read lock the input")?;
    let s: &[u8] = &guard;
    Ok(offset(s))
}

#[allow(unsafe_code)]
#[no_mangle]
pub extern "C" fn offset_i() -> *const u8 {
    _offset_i().ok().unwrap_or_else(std::ptr::null)
}

#[allow(unsafe_code)]
#[no_mangle]
pub extern "C" fn offset_o() -> *const u8 {
    _offset_o().ok().unwrap_or_else(std::ptr::null)
}

#[cfg(feature = "gz_flate2")]
pub fn convert(
    i: &[u8],
    out: &mut Vec<u8>,
    l: CompressionLevel,
    limit: u64,
) -> Result<usize, &'static str> {
    let enc = crate::encode2gzip::fl2::Encoder { level: l, limit };
    enc.encode2gzip(i, out)?;
    Ok(out.len())
}

pub fn _converter(l: CompressionLevel, limit: u64) -> Result<usize, &'static str> {
    let iguard = INPUT
        .try_read()
        .map_err(|_| "unable to read lock the input")?;
    let i: &[u8] = &iguard;
    let mut oguard = OUTPUT
        .try_write()
        .map_err(|_| "unable to write lock the output")?;
    let o: &mut Vec<u8> = &mut oguard;
    convert(i, o, l, limit)
}

macro_rules! create_converter {
    ($fname: ident, $level: expr) => {
        #[allow(unsafe_code)]
        #[no_mangle]
        pub extern "C" fn $fname(limit: u64) -> i32 {
            _converter($level, limit)
                .ok()
                .and_then(|u| u.try_into().ok())
                .unwrap_or(-1)
        }
    };
}

create_converter!(converter_with_limit_fast, CompressionLevel::Fast);
create_converter!(converter_with_limit_none, CompressionLevel::None);
create_converter!(converter_with_limit_best, CompressionLevel::Best);

#[allow(unsafe_code)]
#[no_mangle]
pub extern "C" fn converter_default() -> i32 {
    converter_with_limit_fast(crate::encode2gzip::enc2gz::GZ_SIZE_LIMIT)
}

#[allow(unsafe_code)]
#[no_mangle]
pub extern "C" fn converter() -> i32 {
    converter_default()
}
