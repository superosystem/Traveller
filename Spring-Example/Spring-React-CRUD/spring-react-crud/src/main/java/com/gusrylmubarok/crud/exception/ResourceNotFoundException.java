package com.gusrylmubarok.crud.exception;

import org.springframework.http.HttpStatus;
import org.springframework.web.bind.annotation.ResponseStatus;

/**
 * @author Gusryl Mubarok
 * @email gusrylmubarok@gmail.com
 */

@ResponseStatus(value = HttpStatus.NOT_FOUND)
public class ResourceNotFoundException extends RuntimeException{

    private static final long serialVersionUID = -1873486579019491303L;

    public ResourceNotFoundException(String message) {
        super(message);
    }
}
