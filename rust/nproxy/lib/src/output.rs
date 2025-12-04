use crate::error::NproxyError;
use serde::Serialize;

#[derive(Serialize)]
pub struct ErrorOutput {
    pub error: ErrorDetail,
}

#[derive(Serialize)]
pub struct ErrorDetail {
    pub code: String,
    pub message: String,
}

pub fn print_yaml<T: Serialize>(data: &T) -> Result<(), NproxyError> {
    let yaml = serde_yaml::to_string(data)
        .map_err(|e| NproxyError::ApiError(format!("Failed to serialize: {}", e)))?;
    print!("{}", yaml);
    Ok(())
}

pub fn print_error(err: &NproxyError) {
    let output = ErrorOutput {
        error: ErrorDetail {
            code: err.code().to_string(),
            message: err.to_string(),
        },
    };
    if let Ok(yaml) = serde_yaml::to_string(&output) {
        eprint!("{}", yaml);
    } else {
        eprintln!("error: {}", err);
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[derive(Serialize)]
    struct TestStruct {
        name: String,
        value: i32,
    }

    #[test]
    fn test_yaml_serialization_format() {
        let data = TestStruct { name: "test".into(), value: 42 };
        let yaml = serde_yaml::to_string(&data).unwrap();
        assert!(yaml.contains("name: test"));
        assert!(yaml.contains("value: 42"));
    }

    #[test]
    fn test_error_output_format() {
        let err = NproxyError::ConfigError("test message".into());
        let output = ErrorOutput {
            error: ErrorDetail {
                code: err.code().to_string(),
                message: err.to_string(),
            },
        };
        let yaml = serde_yaml::to_string(&output).unwrap();
        assert!(yaml.contains("code: CONFIG_ERROR"));
        assert!(yaml.contains("message:"));
    }
}
