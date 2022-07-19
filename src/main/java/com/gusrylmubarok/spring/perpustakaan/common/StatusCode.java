package com.gusrylmubarok.spring.perpustakaan.common;

public interface StatusCode {
    String SAVE_SUCCESS = "0010";
    String SAVE_FAILED = "0011";
    String UPDATE_SUCCESS = "0012";
    String UPDATE_FAILED = "0013";
    String DELETE_SUCCESS = "0014";
    String DELETE_FAILED = "0015";

    String OPERATION_SUCCESS = "1112";
    String OPERATION_FAILED = "1113";

    String PASSWORD_OR_USER_NOT_REGISTERED= "0111";
    String LOGIN_SUCCESS = "0112";
    String LOGIN_FAILED = "0113";

}
