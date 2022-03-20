package com.gusrylmubarok.springbasic.event;

import com.gusrylmubarok.springbasic.data.User;
import lombok.Getter;
import org.springframework.context.ApplicationEvent;

/**
 * @author Gusryl Mubarok
 * @email gusrylmubarok@gmail.com
 */

public class LoginSuccessEvent extends ApplicationEvent {

    @Getter
    private final User user;

    public LoginSuccessEvent(User user){
        super(user);
        this.user = user;
    }

}
