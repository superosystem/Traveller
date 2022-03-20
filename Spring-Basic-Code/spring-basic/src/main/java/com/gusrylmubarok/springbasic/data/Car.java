package com.gusrylmubarok.springbasic.data;

import com.gusrylmubarok.springbasic.aware.IdAware;
import lombok.Getter;
import org.springframework.stereotype.Component;

/**
 * @author Gusryl Mubarok
 * @email gusrylmubarok@gmail.com
 */

@Component
public class Car implements IdAware {

    @Getter
    private String id;

    @Override
    public void setId(String id) {
        this.id = id;
    }
}
